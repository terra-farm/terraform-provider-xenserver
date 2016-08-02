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
package main

import (
	"github.com/hashicorp/terraform/helper/schema"
	"log"
	"fmt"
	"bytes"
	"strings"
	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/mborodin/go-xen-api-client"
)

const (
	vbdSchemaUUID                   = "vdi_uuid"
	vbdSchemaBootable               = "bootable"
	vbdSchemaMode                   = "mode"
)

func readVBDsFromSchema(c *Connection, s []interface{}) ([]*VBDDescriptor, error) {
	vbds := make([]*VBDDescriptor, 0, len(s))

	for _, schm := range s {
		data := schm.(map[string]interface{})

		var vdi *VDIDescriptor = nil
		if id, ok := data[vbdSchemaUUID]; ok {
			vdi = &VDIDescriptor{}
			vdi.UUID = id.(string)
			if err := vdi.Load(c); err != nil {
				return nil, err
			}
		}
		bootable := data[vbdSchemaBootable].(bool)

		var mode xenAPI.VbdMode
		_mode := strings.ToLower(data[vbdSchemaMode].(string))

		if _mode == strings.ToLower(string(xenAPI.VbdModeRO)) {
			mode = xenAPI.VbdModeRO
		} else if _mode == strings.ToLower(string(xenAPI.VbdModeRW)) {
			mode = xenAPI.VbdModeRW
		} else {
			return nil, fmt.Errorf("%q is not valid mode (either RO or RW)", data[vbdSchemaMode].(string))
		}

		vbd := &VBDDescriptor{
			VDI: vdi,
			Bootable: bootable,
			Mode: mode,
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
		vbdSchemaUUID: uuid,
		vbdSchemaBootable: vbd.Bootable,
		vbdSchemaMode: vbd.Mode,
	}
}

func createVBD(c *Connection, vbd *VBDDescriptor) (*VBDDescriptor, error) {
	// FIXME: Should be available to add VBD to running VM with PV drivers installed
	// TODO: Check PV driver status
	if vbd.VM.PowerState == xenAPI.VMPowerStateRunning {
		return nil, fmt.Errorf("VM %q(%q) is in running state!", vbd.VM.Name, vbd.VM.UUID)
	}

	log.Println(fmt.Sprintf("[DEBUG] Creating VBD for VM %q", vbd.VM.Name))

	vbdObject := xenAPI.VBDRecord{
		Type: vbd.Type,
		Mode: vbd.Mode,
		Bootable: vbd.Bootable,
		VM: vbd.VM.VMRef,
		Empty: vbd.VDI == nil,
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

	err = c.client.VBD.Plug(c.session, vbdRef)
	if err != nil {
		return nil, err
	}

	log.Println(fmt.Sprintf("[DEBUG] Plugged VBD %q to VM %q", vbd.UUID, vbd.VM.Name))

	return vbd, nil
}

func vbdHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})
	log.Println("[DEBUG] ", m)
	buf.WriteString(fmt.Sprintf("%s-", m["vdi_uuid"].(string)))
	buf.WriteString(fmt.Sprintf("%t-",m["bootable"].(int)))
	buf.WriteString(fmt.Sprintf("%s-",
		strings.ToLower(m["mode"].(string))))

	return hashcode.String(buf.String())
}

func createVBDs(c *Connection, s []interface{}, vbdType xenAPI.VbdType, vm *VMDescriptor) (err error) {

	var vbds []*VBDDescriptor
	if vbds, err = readVBDsFromSchema(c, s); err != nil {
		return err
	}

	for _, vbd := range vbds {
		vbd.Type = vbdType
		vbd.VM = vm

		if _, err = createVBD(c, vbd); err != nil {
			return err
		}
	}

	return nil
}

func resourceVBD() *schema.Resource {
	return &schema.Resource{

		Schema: map[string]*schema.Schema{
			vbdSchemaUUID: &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			vbdSchemaBootable : &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default: false,
			},
			vbdSchemaMode     : &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}