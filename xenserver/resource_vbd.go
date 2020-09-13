/*
 * The MIT License (MIT)
 * Copyright (c) 2016 Maksym Borodin <borodin.maksym@gmail.com>
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated
 * documentation files (the "Software"), to deal in the Software without restriction, including without limitation
 * the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software,
 * and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all copies or substantial portions
 * of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO
 * THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL
 * THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF
 * CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS
 * IN THE SOFTWARE.
 */
package xenserver

import (
	"bytes"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/hashcode"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	xenapi "github.com/terra-farm/go-xen-api-client"
)

const (
	vbdSchemaVdiUUID        = "vdi_uuid"
	vbdSchemaBootable       = "bootable"
	vbdSchemaMode           = "mode"
	vbdSchemaUserDevice     = "user_device"
	vbdSchemaTemplateDevice = "is_from_template"
)

func queryTemplateVBDs(c *Connection, vm *VMDescriptor) (vbds []*VBDDescriptor, err error) {
	vbds = make([]*VBDDescriptor, 0)
	var vmVBDRefs []xenapi.VBDRef
	if vmVBDRefs, err = c.client.VM.GetVBDs(c.session, vm.VMRef); err != nil {
		return
	}

	for _, vmVBDRef := range vmVBDRefs {
		vbd := &VBDDescriptor{
			VBDRef: vmVBDRef,
		}

		if err = vbd.Query(c); err != nil {
			return nil, err
		}

		if vbd.IsTemplateDevice {
			log.Printf("[DEBUG] VBD %s (type = %s) comes from template", vbd.UUID, vbd.Type)
			vbds = append(vbds, vbd)
		}
	}

	log.Printf("[DEBUG] Got %d template vdbs", len(vbds))

	return vbds, nil
}

func readTemplateVBDsToSchema(c *Connection, vm *VMDescriptor, s []interface{}, vbdType xenapi.VbdType) error {
	var vmVBDRefs []xenapi.VBDRef
	var err error
	if vmVBDRefs, err = c.client.VM.GetVBDs(c.session, vm.VMRef); err != nil {
		return err
	}

	for _, vmVBDRef := range vmVBDRefs {
		vbd := &VBDDescriptor{
			VBDRef: vmVBDRef,
		}

		if vbd.Query(c) != nil {
			return err
		}

		// Skip unrelevant VBDs
		if vbdType != vbd.Type {
			continue
		}

		found := false
		for _, schm := range s {
			data := schm.(map[string]interface{})
			userDevice := data[vbdSchemaUserDevice].(string)
			isTemplateDevice := data[vbdSchemaTemplateDevice].(bool)

			if isTemplateDevice && userDevice == vbd.UserDevice {
				found = true

				vbd.IsTemplateDevice = isTemplateDevice

				if err = vbd.Commit(c); err != nil {
					return err
				}

				data[vbdSchemaUserDevice] = vbd.UserDevice
				data[vbdSchemaVdiUUID] = vbd.VDI.UUID
				data[vbdSchemaBootable] = vbd.Bootable
				data[vbdSchemaMode] = vbd.Mode
				data[vbdSchemaTemplateDevice] = isTemplateDevice

				break
			}

		}

		if !found {
			return fmt.Errorf("template VBD %s is not referenced", vbd.UUID)
		}
	}

	return nil
}

func destroyTemplateVDIs(c *Connection, vbds []*VBDDescriptor) (err error) {
	log.Println("[DEBUG] Destroying vbds")
	for _, vbd := range vbds {

		// Only relevant to HDDs
		if vbd.Type != xenapi.VbdTypeDisk {
			continue
		}

		log.Println("[DEBUG] Destroy vbd ", vbd.UUID)
		if err = c.client.VDI.Destroy(c.session, vbd.VDI.VDIRef); err != nil {
			return err
		}
	}
	return nil
}

func readVBDFromSchema(c *Connection, s map[string]interface{}) (*VBDDescriptor, error) {
	// In API it is called user_device, but in terraform provider it is called template device
	// to emphasise that it is used to map VBD from template
	userDevice := s[vbdSchemaUserDevice].(string)

	var vdi *VDIDescriptor = nil

	if id, ok := s[vbdSchemaVdiUUID]; ok {
		log.Println("[DEBUG] Try load VDI ", id)
		vdi = &VDIDescriptor{}
		vdi.UUID = id.(string)
		if err := vdi.Load(c); err != nil {
			return nil, err
		}
	}
	bootable := s[vbdSchemaBootable].(bool)

	var mode xenapi.VbdMode
	_mode := strings.ToLower(s[vbdSchemaMode].(string))

	if _mode == strings.ToLower(string(xenapi.VbdModeRO)) {
		mode = xenapi.VbdModeRO
	} else if _mode == strings.ToLower(string(xenapi.VbdModeRW)) {
		mode = xenapi.VbdModeRW
	} else {
		return nil, fmt.Errorf("%q is not valid mode (either RO or RW)", s[vbdSchemaMode].(string))
	}

	vbd := &VBDDescriptor{
		VDI:        vdi,
		Bootable:   bootable,
		Mode:       mode,
		UserDevice: userDevice,
	}

	return vbd, nil
}

func readVBDsFromSchema(c *Connection, s []interface{}) ([]*VBDDescriptor, error) {
	vbds := make([]*VBDDescriptor, 0, len(s))

	for _, schm := range s {
		data := schm.(map[string]interface{})

		var vbd *VBDDescriptor
		var err error
		if vbd, err = readVBDFromSchema(c, data); err != nil {
			return nil, err
		}
		vbds = append(vbds, vbd)
	}

	return vbds, nil
}

func fillVBDSchema(vbd VBDDescriptor) map[string]interface{} {
	uuid := ""
	if vbd.VDI != nil {
		uuid = vbd.VDI.UUID
	}
	return map[string]interface{}{
		vbdSchemaVdiUUID:        uuid,
		vbdSchemaBootable:       vbd.Bootable,
		vbdSchemaMode:           vbd.Mode,
		vbdSchemaUserDevice:     vbd.UserDevice,
		vbdSchemaTemplateDevice: vbd.IsTemplateDevice,
	}
}

func readVBDs(c *Connection, vm *VMDescriptor) ([]map[string]interface{}, []map[string]interface{}, error) {
	vmVBDs, err := c.client.VM.GetVBDs(c.session, vm.VMRef)
	if err != nil {
		return nil, nil, err
	}

	hdd := make([]map[string]interface{}, 0, len(vmVBDs))
	cdrom := make([]map[string]interface{}, 0, len(vmVBDs))
	log.Println(fmt.Sprintf("[DEBUG] Got %d VBDs", len(vmVBDs)))

	for _, _vbd := range vmVBDs {
		vbd := VBDDescriptor{
			VBDRef: _vbd,
		}

		if err := vbd.Query(c); err != nil {
			return nil, nil, err
		}

		log.Println("[DEBUG] Found VBD", vbd.UUID)
		vbdData := fillVBDSchema(vbd)
		log.Println("[DEBUG] VBD: ", vbdData)
		log.Println("[DEBUG] VBD Type: ", vbd.Type)

		switch vbd.Type {
		case xenapi.VbdTypeCD:
			cdrom = append(cdrom, vbdData)
			break
		case xenapi.VbdTypeDisk:
			hdd = append(hdd, vbdData)
		default:
			return nil, nil, fmt.Errorf("Unsupported VBD type %q", string(vbd.Type))
		}
	}

	return hdd, cdrom, nil
}

func setSchemaVBDs(c *Connection, vm *VMDescriptor, d *schema.ResourceData) error {
	var err error
	var hdd []map[string]interface{}
	var cdrom []map[string]interface{}
	if hdd, cdrom, err = readVBDs(c, vm); err != nil {
		log.Println("[ERROR] ", err)
		return err
	}

	log.Println("[DEBUG] Found ", len(cdrom), " CDs and ", len(hdd), " HDDs")
	err = d.Set(vmSchemaHardDrive, hdd)
	if err != nil {
		log.Println("[ERROR] ", err)
		return err
	}
	err = d.Set(vmSchemaCdRom, cdrom)
	if err != nil {
		log.Println("[ERROR] ", err)
		return err
	}

	return nil
}

func createVBD(c *Connection, vbd *VBDDescriptor) (*VBDDescriptor, error) {
	log.Println(fmt.Sprintf("[DEBUG] Creating VBD for VM %q", vbd.VM.Name))

	vbdObject := xenapi.VBDRecord{
		Type:       vbd.Type,
		Mode:       vbd.Mode,
		Bootable:   vbd.Bootable,
		VM:         vbd.VM.VMRef,
		Empty:      vbd.VDI == nil,
		Userdevice: vbd.UserDevice,
	}

	if devices, err := c.client.VM.GetAllowedVBDDevices(c.session, vbd.VM.VMRef); err == nil {
		if len(devices) == 0 {
			return nil, fmt.Errorf("No available devices to attach to")
		}
		vbdObject.Userdevice = devices[0]
		log.Println("[DEBUG] Selected device for VBD: ", vbdObject.Userdevice)
	} else {
		return nil, err
	}

	if vbd.VDI != nil {
		vbdObject.VDI = vbd.VDI.VDIRef
	}

	vbdRef, err := c.client.VBD.Create(c.session, vbdObject)
	if err != nil {
		return nil, err
	}

	log.Println(fmt.Sprintf("[DEBUG] Created VBD"))

	vbd.VBDRef = vbdRef
	err = vbd.Query(c)
	if err != nil {
		return nil, err
	}

	log.Println(fmt.Sprintf("[DEBUG] VBD  UUID %q", vbd.UUID))

	if vbd.VM.PowerState == xenapi.VMPowerStateRunning {
		err = c.client.VBD.Plug(c.session, vbdRef)
		if err != nil {
			return nil, err
		}

		log.Println(fmt.Sprintf("[DEBUG] Plugged VBD %q to VM %q", vbd.UUID, vbd.VM.Name))
	}

	return vbd, nil
}

func vbdHash(v interface{}) int {
	m := v.(map[string]interface{})
	var buf bytes.Buffer
	var count int = 0
	var b int

	userDevice := m[vbdSchemaUserDevice].(string)
	isTemplateDevice := m[vbdSchemaTemplateDevice].(bool)
	mode := m[vbdSchemaMode].(string)
	bootable := m[vbdSchemaBootable].(bool)
	vdiUUID := m[vbdSchemaVdiUUID].(string)

	log.Println("[DEBUG] Calculating hash for ", v)

	if isTemplateDevice || userDevice != "" {
		b, _ = buf.WriteString(fmt.Sprintf("%s", userDevice))
		count += b
	}

	if !isTemplateDevice {
		b, _ = buf.WriteString(fmt.Sprintf("-%s", vdiUUID))
		count += b

		if mode != "" {
			b, _ = buf.WriteString(fmt.Sprintf("%-s", mode))
			count += b
		}

		b, _ = buf.WriteString(fmt.Sprintf("-%t", bootable))
		count += b
	}
	log.Println("Consumed total ", count, " bytes to generate hash")
	log.Println("String for hash: ", buf.String())

	return hashcode.String(buf.String())
}

func createVBDs(c *Connection, s []interface{}, vbdType xenapi.VbdType, vm *VMDescriptor) (err error) {

	if err := readTemplateVBDsToSchema(c, vm, s, vbdType); err != nil {
		return err
	}

	log.Println("[DEBUG] Creating ", len(s), " VBDS of type ", vbdType)

	for _, schm := range s {
		data := schm.(map[string]interface{})

		if data[vbdSchemaUserDevice] != "" {
			continue
		}

		var vbd *VBDDescriptor
		var err error
		if vbd, err = readVBDFromSchema(c, data); err != nil {
			return err
		}

		// User Device is used to specify existing VBDs
		if vbd.UserDevice != "" {
			continue
		}

		vbd.Type = vbdType
		vbd.VM = vm

		if vbdType == xenapi.VbdTypeCD {
			vbd.Mode = xenapi.VbdModeRO
		}

		if vbd, err = createVBD(c, vbd); err != nil {
			return err
		}

		data[vbdSchemaUserDevice] = vbd.UserDevice
		data[vbdSchemaVdiUUID] = vbd.VDI.UUID
		data[vbdSchemaBootable] = vbd.Bootable
		data[vbdSchemaMode] = vbd.Mode
	}

	return nil
}

func resourceVBD() *schema.Resource {
	return &schema.Resource{

		Schema: map[string]*schema.Schema{
			vbdSchemaTemplateDevice: &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			vbdSchemaVdiUUID: &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			vbdSchemaUserDevice: &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			vbdSchemaBootable: &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			vbdSchemaMode: &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}
