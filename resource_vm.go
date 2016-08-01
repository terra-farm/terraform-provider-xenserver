/*
 * The MIT License (MIT)
 * Copyright (c) 2016 Michael Franz Aigner <maigner@updox.com>
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
	"fmt"
	"github.com/amfranz/go-xen-api-client"
	"github.com/hashicorp/terraform/helper/schema"
	"strings"
)

const (
	vmSchemaNameLabel                 = "name_label"
	vmSchemaBaseTemplateName          = "base_template_name"
	vmSchemaMemory                    = "mem"
	vmSchemaStaticMemoryMin           = "static_mem_min"
	vmSchemaStaticMemoryMax           = "static_mem_max"
	vmSchemaDynamicMemoryMin          = "dynamic_mem_min"
	vmSchemaDynamicMemoryMax          = "dynamic_mem_max"
	vmSchemaBootOrder                 = "boot_order"
	vmSchemaNetworkInterfaces         = "network_interface"
	vmSchemaHardDrives                = "hard_drives"
	vmSchemaBootParameters            = "boot_parameters"
	vmSchemaInstallationMediaType     = "installation_media_type"
	vmSchemaInstallationMediaLocation = "installation_media_location"
	vmSchemaVcpus                     = "vcpus"
	vmSchemaCoresPerSocket            = "cores_per_socket"
	vmSchemaXenstoreData              = "xenstore_data"
)

const xenstoreVMDataPrefix = "vm-data/"

func resourceVM() *schema.Resource {
	return &schema.Resource{
		Create: resourceVMCreate,
		Read:   resourceVMRead,
		Update: resourceVMUpdate,
		Delete: resourceVMDelete,
		Exists: resourceVMExists,

		Schema: map[string]*schema.Schema{
			vmSchemaNameLabel: &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			vmSchemaBaseTemplateName: &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			vmSchemaXenstoreData: &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
			},

			vmSchemaMemory: &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},

			vmSchemaStaticMemoryMin: &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},

			vmSchemaStaticMemoryMax: &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},

			vmSchemaDynamicMemoryMin: &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},

			vmSchemaDynamicMemoryMax: &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},

			vmSchemaBootOrder: &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default: "dc",
			},

			vmSchemaNetworkInterfaces: &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				Elem: resourceVIF(),
			},

			vmSchemaHardDrives: &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},

			vmSchemaBootParameters: &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},

			vmSchemaInstallationMediaType: &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},

			vmSchemaInstallationMediaLocation: &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},

			vmSchemaVcpus: &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},

			vmSchemaCoresPerSocket: &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},

		},
	}
}

func filterVMTemplates(c *Connection, vms []xenAPI.VMRef) ([]xenAPI.VMRef, error) {
	var templates []xenAPI.VMRef
	for _, vm := range vms {
		isATemplate, err := c.client.VM.GetIsATemplate(c.session, vm)
		if err != nil {
			return templates, err
		}
		if isATemplate {
			templates = append(templates, vm)
		}
	}
	return templates, nil
}

func resourceVMCreate(d *schema.ResourceData, m interface{}) error {
	c := m.(*Connection)

	dBaseTemplateName := d.Get(vmSchemaBaseTemplateName).(string)

	xenBaseTemplates, err := c.client.VM.GetByNameLabel(c.session, dBaseTemplateName)
	if err != nil {
		return err
	}

	xenBaseTemplates, err = filterVMTemplates(c, xenBaseTemplates)
	if err != nil {
		return err
	}

	if len(xenBaseTemplates) == 0 {
		return fmt.Errorf("No VM template with label %q has been found", dBaseTemplateName)
	}

	if len(xenBaseTemplates) > 1 {
		return fmt.Errorf("More than one VM template with label %q has been found", dBaseTemplateName)
	}

	xenBaseTemplate := xenBaseTemplates[0]

	dNameLabel := d.Get(vmSchemaNameLabel).(string)

	xenVM, err := c.client.VM.Clone(c.session, xenBaseTemplate, dNameLabel)
	if err != nil {
		return err
	}

	vm := &VMDescriptor{
		VMRef: xenVM,
	}

	if err = vm.Query(c); err != nil {
		return err
	}

	// Memory configuration
	mem := d.Get(vmSchemaMemory)
	vm.StaticMemory = Range{
		Min: mem.(int),
		Max: mem.(int),
	}
	vm.DynamicMemory = Range{
		Min: mem.(int),
		Max: mem.(int),
	}

	mem, ok := d.GetOk(vmSchemaStaticMemoryMin)
	if ok {
		vm.StaticMemory.Min = mem.(int)
	}

	mem, ok = d.GetOk(vmSchemaStaticMemoryMax)
	if ok {
		vm.StaticMemory.Max = mem.(int)
	}

	mem, ok = d.GetOk(vmSchemaDynamicMemoryMin)
	if ok {
		vm.DynamicMemory.Min = mem.(int)
	}

	mem, ok = d.GetOk(vmSchemaDynamicMemoryMax)
	if ok {
		vm.DynamicMemory.Max = mem.(int)
	}

	if err=vm.UpdateMemory(c); err != nil {
		return err
	}

	// Set VCPUs number
	vm.VCPUCount = d.Get(vmSchemaVcpus).(int)
	if err = vm.UpdateVCPUs(c); err != nil {
		return err
	}

	d.SetId(vm.UUID)

	// TODO: Refactor
	dXenstoreDataRaw, ok := d.GetOk(vmSchemaXenstoreData)
	if ok {
		dXenstoreData := make(map[string]string)
		for key, value := range dXenstoreDataRaw.(map[string]interface{}) {
			dXenstoreData[xenstoreVMDataPrefix+key] = value.(string)
		}

		err = c.client.VM.SetXenstoreData(c.session, xenVM, dXenstoreData)
		if err != nil {
			return err
		}
	}

	var vifs []*VIFDescriptor

	if vifs, err = readVIFsFromSchema(c, d.Get(vmSchemaNetworkInterfaces).(*schema.Set).List()); err != nil {
		return err
	}

	for _, vif := range vifs {
		vif.VM = vm
		if vif, err = createVIF(c, vif); err != nil {
			return nil
		}
	}

	err = c.client.VM.Provision(c.session, xenVM)
	if err != nil {
		return err
	}

	//err = c.client.VM.Start(c.session, xenVM, false, false)
	//if err != nil {
	//	return err
	//}

	return nil
}

func resourceVMRead(d *schema.ResourceData, m interface{}) error {
	c := m.(*Connection)

	xenVM, err := c.client.VM.GetByUUID(c.session, d.Id())
	if err != nil {
		if xenErr, ok := err.(*xenAPI.Error); ok {
			if xenErr.Code() == xenAPI.ERR_UUID_INVALID {
				d.SetId("")
				return nil
			}
		}

		return err
	}

	vmRecord, err := c.client.VM.GetRecord(c.session, xenVM)
	if err != nil {
		return err
	}

	err = d.Set(vmSchemaNameLabel, vmRecord.NameLabel)
	if err != nil {
		return err
	}

	vmBaseTemplateName, ok := vmRecord.OtherConfig["base_template_name"]
	if ok {
		err = d.Set(vmSchemaBaseTemplateName, vmBaseTemplateName)
		if err != nil {
			return err
		}
	}

	vmXenstoreData := make(map[string]string)
	for key, value := range vmRecord.XenstoreData {
		if strings.HasPrefix(key, xenstoreVMDataPrefix) {
			vmXenstoreData[key[len(xenstoreVMDataPrefix):]] = value
		}
	}

	err = d.Set(vmSchemaXenstoreData, vmXenstoreData)
	if err != nil {
		return err
	}

	return nil
}

func resourceVMUpdate(d *schema.ResourceData, m interface{}) error {
	c := m.(*Connection)

	vm := &VMDescriptor{
		UUID: d.Id(),
	}
	if err := vm.Load(c); err != nil {
		if xenErr, ok := err.(*xenAPI.Error); ok {
			if xenErr.Code() == xenAPI.ERR_UUID_INVALID {
				d.SetId("")
				return nil
			}
		}

		return err
	}

	d.Partial(true)

	if d.HasChange(vmSchemaNameLabel) {
		_, _dNameLabel := d.GetChange(vmSchemaNameLabel)
		dNameLabel := _dNameLabel.(string)
		if err := c.client.VM.SetNameLabel(c.session, vm.VMRef, dNameLabel); err != nil {
			return err
		}

		d.SetPartial(vmSchemaNameLabel)
	}

	updatedFields := make([]string, 0, 5)
	updateMemory := false
	if d.HasChange(vmSchemaMemory) {
		mem := d.Get(vmSchemaMemory).(int)
		vm.StaticMemory = Range{
			Max: mem,
			Min: mem,
		}
		vm.DynamicMemory = Range{
			Max: mem,
			Min: mem,
		}
		updateMemory = true
		updatedFields = append(updatedFields, vmSchemaMemory)
	}

	if d.HasChange(vmSchemaStaticMemoryMax) {
		mem := d.Get(vmSchemaStaticMemoryMax).(int)
		vm.StaticMemory.Max = mem
		updateMemory = true
		updatedFields = append(updatedFields, vmSchemaStaticMemoryMax)
	}

	if d.HasChange(vmSchemaStaticMemoryMin) {
		mem := d.Get(vmSchemaStaticMemoryMin).(int)
		vm.StaticMemory.Min = mem
		updateMemory = true
		updatedFields = append(updatedFields, vmSchemaStaticMemoryMin)
	}

	if d.HasChange(vmSchemaDynamicMemoryMax) {
		mem := d.Get(vmSchemaDynamicMemoryMax).(int)
		vm.DynamicMemory.Max = mem
		updateMemory = true
		updatedFields = append(updatedFields, vmSchemaDynamicMemoryMax)
	}

	if d.HasChange(vmSchemaDynamicMemoryMin) {
		mem := d.Get(vmSchemaDynamicMemoryMin).(int)
		vm.DynamicMemory.Min = mem
		updateMemory = true
		updatedFields = append(updatedFields, vmSchemaDynamicMemoryMin)
	}

	if updateMemory {
		if err := vm.UpdateMemory(c); err != nil {
			return err
		}

		for _, f := range updatedFields {
			d.SetPartial(f)
		}
	}

	if d.HasChange(vmSchemaNetworkInterfaces) {
		o, n := d.GetChange(vmSchemaNetworkInterfaces)

		os := o.(*schema.Set)
		ns := n.(*schema.Set)

		var err error
		var remove []*VIFDescriptor
		if remove, err = readVIFsFromSchema(c, os.Difference(ns).List()); err == nil {
			return err
		}

		if len(remove) > 0 {

			var vmVifs []*VIFDescriptor
			if _vmVifs, err := c.client.VM.GetVIFs(c.session, vm.VMRef); err == nil {
				for _, _vif := range _vmVifs {
					vif := &VIFDescriptor{
						VIFRef: _vif,
					}

					if err := vif.Query(c); err != nil {
						return err
					}
					vmVifs = append(vmVifs, vif)
				}
			} else {
				return err
			}

			for _, vif := range remove {
				var vifToRemove *VIFDescriptor
				for _, candidate := range vmVifs {
					if candidate.Network.UUID == vif.Network.UUID && candidate.DeviceOrder == vif.DeviceOrder {
						vifToRemove = candidate
						break
					}
				}
				if vifToRemove != nil {
					if err := c.client.VIF.Destroy(c.session, vifToRemove.VIFRef); err != nil {
						return err
					}
				}
			}
		}

		var create []*VIFDescriptor
		if create, err = readVIFsFromSchema(c, ns.Difference(os).List()); err == nil {
			return err
		}

		if len(create) > 0 {
			for _, vif := range create {
				vif.VM = vm
				if _, err := createVIF(c, vif); err != nil {
					return nil
				}
			}
		}

	}

	dXenstoreDataRaw, ok := d.GetOk(vmSchemaXenstoreData)
	if ok {
		dXenstoreData := make(map[string]string)
		for key, value := range dXenstoreDataRaw.(map[string]interface{}) {
			dXenstoreData[xenstoreVMDataPrefix+key] = value.(string)
		}

		if err := c.client.VM.SetXenstoreData(c.session, vm.VMRef, dXenstoreData); err != nil {
			return err
		}

		d.SetPartial(vmSchemaXenstoreData)
	}

	d.Partial(false)

	return resourceVMRead(d, m)
}

func resourceVMDelete(d *schema.ResourceData, m interface{}) error {
	c := m.(*Connection)

	vm := VMDescriptor{
		UUID: d.Id(),
	}
	if err := vm.Load(c); err != nil {
		if xenErr, ok := err.(*xenAPI.Error); ok {
			if xenErr.Code() == xenAPI.ERR_UUID_INVALID {
				d.SetId("")
				return nil
			}
		}

		return err
	}

	if vm.PowerState == xenAPI.VMPowerStateRunning {
		if err := c.client.VM.HardShutdown(c.session, vm.VMRef);  err != nil {
			return err
		}
	}

	if err := c.client.VM.Destroy(c.session, vm.VMRef); err != nil {
		return err
	}

	d.SetId("")
	return nil
}

func resourceVMExists(d *schema.ResourceData, m interface{}) (bool, error) {
	c := m.(*Connection)

	_, err := c.client.VM.GetByUUID(c.session, d.Id())
	if err != nil {
		if xenErr, ok := err.(*xenAPI.Error); ok {
			if xenErr.Code() == xenAPI.ERR_UUID_INVALID {
				return false, nil
			}
		}

		return false, err
	}

	return true, nil
}
