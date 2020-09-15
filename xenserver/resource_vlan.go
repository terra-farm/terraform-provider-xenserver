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
	pifSchemaUUID         = "uuid"
	vlanSchemaUUID        = "uuid"
	vlanSchemaTag         = "tag"
	vlanSchemaPIF         = "pif"
	vlanSchemaOtherConfig = "other_config"
	vlanSchemaNetwork     = "network"
)

func resourceVLAN() *schema.Resource {
	return &schema.Resource{
		Create: resourceVLANCreate,
		Read:   resourceVLANRead,
		Update: resourceVLANUpdate,
		Delete: resourceVLANDelete,
		Exists: resourceVLANExists,

		Schema: map[string]*schema.Schema{
			vlanSchemaTag: &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},

			vlanSchemaPIF: &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			vlanSchemaNetwork: &schema.Schema{
				ForceNew: true,
				Type:     schema.TypeString,
				Required: true,
			},

			vlanSchemaOtherConfig: &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
			},
		},
	}
}

func resourceVLANCreate(d *schema.ResourceData, m interface{}) error {
	c := m.(*Connection)

	pif := PIFDescriptor{
		UUID: d.Get(vlanSchemaPIF).(string),
	}

	if err := pif.Load(c); err != nil {
		return err
	}

	network := NetworkDescriptor{
		UUID: d.Get(vlanSchemaNetwork).(string),
	}

	if err := network.Load(c); err != nil {
		return err
	}

	tag := d.Get(vlanSchemaTag).(int)

	if vlanRef, err := c.client.VLAN.Create(c.session, pif.PIFRef, tag, network.NetworkRef); err == nil {
		log.Println("VLAN Created")
		vlan := &VLANDescriptor{
			VLANRef: vlanRef,
		}

		if err := vlan.Query(c); err != nil {
			return err
		}
		log.Println("UUID is ", vlan.UUID)
		d.SetId(vlan.UUID)

		if _otherConfig, ok := d.GetOk(vlanSchemaOtherConfig); ok {
			otherConfig := _otherConfig.(map[string]string)
			for k, v := range otherConfig {
				if err := c.client.VLAN.AddToOtherConfig(c.session, vlan.VLANRef, k, v); err != nil {
					return err
				}
			}
		}
	} else {
		log.Println("VLAN not created!")
		return err
	}

	return nil
}

func resourceVLANRead(d *schema.ResourceData, m interface{}) error {
	c := m.(*Connection)

	vlan := &VLANDescriptor{
		UUID: d.Id(),
	}

	if err := vlan.Load(c); err != nil {
		return err
	}

	d.SetId(vlan.UUID)
	if err := d.Set(vlanSchemaTag, vlan.Tag); err != nil {
		return err
	}

	if err := d.Set(vlanSchemaOtherConfig, vlan.OtherConfig); err != nil {
		return err
	}

	if err := d.Set(vlanSchemaPIF, vlan.UntaggedPIF); err != nil {
		return err
	}

	return nil
}
func resourceVLANUpdate(d *schema.ResourceData, m interface{}) error {
	c := m.(*Connection)

	vlan := &VLANDescriptor{
		UUID: d.Id(),
	}

	if err := vlan.Load(c); err != nil {
		return err
	}

	if d.HasChange(vlanSchemaOtherConfig) {
		_, n := d.GetChange(vlanSchemaOtherConfig)
		otherConfig := make(map[string]string)

		for k, v := range n.(map[string]interface{}) {
			otherConfig[k] = v.(string)
		}

		if err := c.client.VLAN.SetOtherConfig(c.session, vlan.VLANRef, otherConfig); err != nil {
			return err
		}

		d.SetPartial(vlanSchemaOtherConfig)
	}

	return nil
}
func resourceVLANDelete(d *schema.ResourceData, m interface{}) error {
	c := m.(*Connection)

	vlan := &VLANDescriptor{
		UUID: d.Id(),
	}

	if err := vlan.Load(c); err != nil {
		return err
	}

	if err := c.client.VLAN.Destroy(c.session, vlan.VLANRef); err != nil {
		return err
	}

	return nil
}
func resourceVLANExists(d *schema.ResourceData, m interface{}) (bool, error) {
	c := m.(*Connection)

	vlan := &VLANDescriptor{
		UUID: d.Id(),
	}

	if err := vlan.Load(c); err != nil {
		if xenErr, ok := err.(*xenapi.Error); ok {
			if xenErr.Code() == xenapi.ERR_UUID_INVALID {
				return false, nil
			}
		}

		return false, err
	}

	return true, nil
}
