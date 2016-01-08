package main

import (
	"fmt"
	"github.com/amfranz/go-xen-api-client"
	"github.com/hashicorp/terraform/helper/schema"
	"strings"
)

const (
	vmSchemaNameLabel        = "name_label"
	vmSchemaBaseTemplateName = "base_template_name"
	vmSchemaXenstoreData     = "xenstore_data"
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
		},
	}
}

func filterVMTemplates(c *Connection, vms []xenAPI.VMRef) ([]xenAPI.VMRef, error) {
	var templates []xenAPI.VMRef
	for _, vm := range vms {
		isATemplate, err := c.client.VM().GetIsATemplate(c.session, vm)
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

	xenBaseTemplates, err := c.client.VM().GetByNameLabel(c.session, dBaseTemplateName)
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

	xenVM, err := c.client.VM().Clone(c.session, xenBaseTemplate, dNameLabel)
	if err != nil {
		return err
	}

	xenUUID, err := c.client.VM().GetUUID(c.session, xenVM)
	if err != nil {
		return err
	}

	d.SetId(xenUUID)

	dXenstoreDataRaw, ok := d.GetOk(vmSchemaXenstoreData)
	if ok {
		dXenstoreData := make(map[string]string)
		for key, value := range dXenstoreDataRaw.(map[string]interface{}) {
			dXenstoreData[xenstoreVMDataPrefix+key] = value.(string)
		}

		err = c.client.VM().SetXenstoreData(c.session, xenVM, dXenstoreData)
		if err != nil {
			return err
		}
	}

	err = c.client.VM().Provision(c.session, xenVM)
	if err != nil {
		return err
	}

	return nil
}

func resourceVMRead(d *schema.ResourceData, m interface{}) error {
	c := m.(*Connection)

	xenVM, err := c.client.VM().GetByUUID(c.session, d.Id())
	if err != nil {
		if xenErr, ok := err.(*xenAPI.Error); ok {
			if xenErr.Code() == xenAPI.ERR_UUID_INVALID {
				d.SetId("")
				return nil
			}
		}

		return err
	}

	vmRecord, err := c.client.VM().GetRecord(c.session, xenVM)
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

	xenVM, err := c.client.VM().GetByUUID(c.session, d.Id())
	if err != nil {
		if xenErr, ok := err.(*xenAPI.Error); ok {
			if xenErr.Code() == xenAPI.ERR_UUID_INVALID {
				d.SetId("")
				return nil
			}
		}

		return err
	}

	vmRecord, err := c.client.VM().GetRecord(c.session, xenVM)
	if err != nil {
		return err
	}

	d.Partial(true)

	dNameLabel := d.Get(vmSchemaNameLabel).(string)
	if vmRecord.NameLabel != dNameLabel {
		err := c.client.VM().SetNameLabel(c.session, xenVM, dNameLabel)
		if err != nil {
			return err
		}

		d.SetPartial(vmSchemaNameLabel)
	}

	dXenstoreDataRaw, ok := d.GetOk(vmSchemaXenstoreData)
	if ok {
		dXenstoreData := make(map[string]string)
		for key, value := range dXenstoreDataRaw.(map[string]interface{}) {
			dXenstoreData[xenstoreVMDataPrefix+key] = value.(string)
		}

		err = c.client.VM().SetXenstoreData(c.session, xenVM, dXenstoreData)
		if err != nil {
			return err
		}

		d.SetPartial(vmSchemaXenstoreData)
	}

	d.Partial(false)

	return resourceVMRead(d, m)
}

func resourceVMDelete(d *schema.ResourceData, m interface{}) error {
	c := m.(*Connection)

	xenVM, err := c.client.VM().GetByUUID(c.session, d.Id())
	if err != nil {
		if xenErr, ok := err.(*xenAPI.Error); ok {
			if xenErr.Code() == xenAPI.ERR_UUID_INVALID {
				d.SetId("")
				return nil
			}
		}

		return err
	}

	powerState, err := c.client.VM().GetPowerState(c.session, xenVM)
	if err != nil {
		return err
	}

	if powerState == xenAPI.VMPowerStateRunning {
		err = c.client.VM().HardShutdown(c.session, xenVM)
		if err != nil {
			return err
		}
	}

	err = c.client.VM().Destroy(c.session, xenVM)
	if err != nil {
		return err
	}

	d.SetId("")
	return nil
}

func resourceVMExists(d *schema.ResourceData, m interface{}) (bool, error) {
	c := m.(*Connection)

	_, err := c.client.VM().GetByUUID(c.session, d.Id())
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
