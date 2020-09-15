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
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	xenapi "github.com/terra-farm/go-xen-api-client"
)

const (
	vdiSchemaUUID   = "sr_uuid"
	vdiSchemaName   = "name_label"
	vdiSchemaShared = "shared"
	vdiSchemaRO     = "read_only"
	vdiSchemaSize   = "size"
)

func resourceVDI() *schema.Resource {
	return &schema.Resource{
		Create: resourceVDICreate,
		Read:   resourceVDIRead,
		Update: resourceVDIUpdate,
		Delete: resourceVDIDelete,
		Exists: resourceVDIExists,

		Schema: map[string]*schema.Schema{
			vdiSchemaUUID: &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			vdiSchemaName: &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			vdiSchemaShared: &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			vdiSchemaRO: &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			vdiSchemaSize: &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
		},
	}
}

func resourceVDICreate(d *schema.ResourceData, m interface{}) error {
	c := m.(*Connection)

	sr := &SRDescriptor{
		UUID: d.Get(vdiSchemaUUID).(string),
	}

	log.Println("Going to create VDI in SR ", sr.UUID)

	if err := sr.Load(c); err != nil {
		log.Println("SR not found!")
		return err
	}

	vdiRecord := xenapi.VDIRecord{
		NameLabel:   d.Get(vdiSchemaName).(string),
		VirtualSize: d.Get(vdiSchemaSize).(int),
		Sharable:    d.Get(vdiSchemaShared).(bool),
		ReadOnly:    d.Get(vdiSchemaRO).(bool),
		SR:          sr.SRRef,
		Type:        xenapi.VdiTypeUser,
	}

	log.Println("Object to send: ", vdiRecord)
	if vdiRef, err := c.client.VDI.Create(c.session, vdiRecord); err == nil {
		log.Println("VDI Created")
		vdi := &VDIDescriptor{
			VDIRef: vdiRef,
		}

		if err := vdi.Query(c); err != nil {
			return err
		}
		log.Println("UUID is ", vdi.UUID)
		d.SetId(vdi.UUID)
	} else {
		log.Println("VDI not created!")
		return err
	}

	return nil
}

func resourceVDIRead(d *schema.ResourceData, m interface{}) error {
	c := m.(*Connection)

	vdi := &VDIDescriptor{
		UUID: d.Id(),
	}

	if err := vdi.Load(c); err != nil {
		return err
	}

	d.SetId(vdi.UUID)
	if err := d.Set(vdiSchemaName, vdi.Name); err != nil {
		return err
	}

	if err := d.Set(vdiSchemaRO, vdi.IsReadOnly); err != nil {
		return err
	}

	if err := d.Set(vdiSchemaShared, vdi.IsShared); err != nil {
		return err
	}

	if err := d.Set(vdiSchemaSize, vdi.Size); err != nil {
		return err
	}

	return nil
}
func resourceVDIUpdate(d *schema.ResourceData, m interface{}) error {
	c := m.(*Connection)

	vdi := &VDIDescriptor{
		UUID: d.Id(),
	}

	if err := vdi.Load(c); err != nil {
		return err
	}

	if d.HasChange(vdiSchemaName) {
		_, n := d.GetChange(vdiSchemaName)

		if err := c.client.VDI.SetNameLabel(c.session, vdi.VDIRef, n.(string)); err != nil {
			return err
		}

		d.SetPartial(vdiSchemaName)
	}

	if d.HasChange(vdiSchemaSize) {
		_, n := d.GetChange(vdiSchemaSize)

		if err := c.client.VDI.SetVirtualSize(c.session, vdi.VDIRef, n.(int)); err != nil {
			return err
		}

		d.SetPartial(vdiSchemaSize)
	}

	if d.HasChange(vdiSchemaShared) {
		_, n := d.GetChange(vdiSchemaShared)

		if err := c.client.VDI.SetSharable(c.session, vdi.VDIRef, n.(bool)); err != nil {
			return err
		}

		d.SetPartial(vdiSchemaShared)
	}

	if d.HasChange(vdiSchemaRO) {
		_, n := d.GetChange(vdiSchemaRO)

		if err := c.client.VDI.SetReadOnly(c.session, vdi.VDIRef, n.(bool)); err != nil {
			return err
		}

		d.SetPartial(vdiSchemaRO)
	}

	return nil
}
func resourceVDIDelete(d *schema.ResourceData, m interface{}) error {
	c := m.(*Connection)

	vdi := &VDIDescriptor{
		UUID: d.Id(),
	}

	if err := vdi.Load(c); err != nil {
		return err
	}

	if err := c.client.VDI.Destroy(c.session, vdi.VDIRef); err != nil {
		return err
	}

	return nil
}
func resourceVDIExists(d *schema.ResourceData, m interface{}) (bool, error) {
	c := m.(*Connection)

	vdi := &VDIDescriptor{
		UUID: d.Id(),
	}

	if err := vdi.Load(c); err != nil {
		if xenErr, ok := err.(*xenapi.Error); ok {
			if xenErr.Code() == xenapi.ERR_UUID_INVALID {
				return false, nil
			}
		}

		return false, err
	}

	return true, nil
}
