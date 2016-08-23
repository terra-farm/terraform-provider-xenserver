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
	"github.com/mborodin/go-xen-api-client"
	"log"
)

const (
	networkSchemaName                 = "name_label"
	networkSchemaDescription          = "description"
	networkSchemaPhysicalInterface    = "physical_interface"
	networkSchemaMTU                  = "mtu"
)

func resourceNetwork() *schema.Resource {
	return &schema.Resource{
		Create: resourceNetworkCreate,
		Read:   resourceNetworkRead,
		Update: resourceNetworkUpdate,
		Delete: resourceNetworkDelete,
		Exists: resourceNetworkExists,

		Schema: map[string]*schema.Schema{
			networkSchemaName  : &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			networkSchemaDescription: &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},

			networkSchemaPhysicalInterface    : &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{Type: schema.TypeString},
				Set: schema.HashString,
			},

			networkSchemaMTU  : &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},

		},
	}
}

func resourceNetworkCreate(d *schema.ResourceData, m interface{}) error {
	c := m.(*Connection)

	networkRecord := xenAPI.NetworkRecord{
		NameLabel: d.Get(networkSchemaName).(string),
		NameDescription: d.Get(networkSchemaDescription).(string),
		MTU: d.Get(networkSchemaMTU).(int),
	}

	if _interfaces, ok := d.GetOk(networkSchemaPhysicalInterface); ok {
		interfaces := _interfaces.(*schema.Set).List()
		pifs := make([]xenAPI.PIFRef, len(interfaces));

		for _, _interface := range interfaces {
			_int := _interface.(string)

			if pif, err := c.client.PIF.GetByUUID(c.session, _int); err == nil {
				pifs = append(pifs, pif)
			} else {
				return err
			}
		}

		networkRecord.PIFs = pifs
	}

	if networkRef, err := c.client.Network.Create(c.session, networkRecord); err == nil {
		log.Println("Network Created")
		network := &NetworkDescriptor{
			NetworkRef: networkRef,
		}

		if err := network.Query(c); err != nil {
			return err
		}
		log.Println("UUID is ", network.UUID)
		d.SetId(network.UUID)
	} else {
		log.Println("Network not created!")
		return err
	}

	return nil
}

func resourceNetworkRead(d *schema.ResourceData, m interface{}) error {
	c := m.(*Connection)

	network := &NetworkDescriptor{
		UUID: d.Id(),
	}

	if err := network.Load(c); err != nil {
		return err
	}

	d.SetId(network.UUID)
	if err := d.Set(networkSchemaName, network.Name); err != nil {
		return err
	}

	if err := d.Set(networkSchemaPhysicalInterface, network.Bridge); err != nil {
		return err
	}

	if err := d.Set(networkSchemaMTU, network.MTU); err != nil {
		return err
	}

	if err := d.Set(networkSchemaDescription, network.Description); err != nil {
		return err
	}

	if pifs, err := c.client.Network.GetPIFs(c.session, network.NetworkRef); err == nil {
		pifUUIDs := &schema.Set{}
		pifUUIDs.F = schema.HashString

		for _, pif := range pifs {
			if pifRecord, err := c.client.PIF.GetRecord(c.session, pif); err == nil {
				pifUUIDs.Add(pifRecord.UUID)
			} else {
				return err
			}
		}
		if err := d.Set(networkSchemaPhysicalInterface, pifUUIDs); err != nil {
			return err
		}

	} else {
		return err
	}

	return nil
}
func resourceNetworkUpdate(d *schema.ResourceData, m interface{}) error {
	c := m.(*Connection)

	network := &NetworkDescriptor{
		UUID: d.Id(),
	}

	if err := network.Load(c); err != nil {
		return err
	}

	if(d.HasChange(networkSchemaName)) {
		_, n := d.GetChange(networkSchemaName)

		if err := c.client.Network.SetNameLabel(c.session, network.NetworkRef, n.(string)); err != nil {
			return err;
		}

		d.SetPartial(networkSchemaName)
	}

	if(d.HasChange(networkSchemaMTU)) {
		_, n := d.GetChange(networkSchemaMTU)

		if err := c.client.Network.SetMTU(c.session, network.NetworkRef, n.(int)); err != nil {
			return err;
		}

		d.SetPartial(networkSchemaMTU)
	}

	if(d.HasChange(networkSchemaDescription)) {
		_, n := d.GetChange(networkSchemaDescription)

		if err := c.client.Network.SetNameDescription(c.session, network.NetworkRef, n.(string)); err != nil {
			return err;
		}

		d.SetPartial(networkSchemaDescription)
	}

	return nil
}
func resourceNetworkDelete(d *schema.ResourceData, m interface{}) error {
	c := m.(*Connection)

	network := &NetworkDescriptor{
		UUID: d.Id(),
	}

	if err := network.Load(c); err != nil {
		return err
	}

	if err := c.client.Network.Destroy(c.session, network.NetworkRef); err != nil {
		return err
	}

	return nil
}
func resourceNetworkExists(d *schema.ResourceData, m interface{}) (bool,error) {
	c := m.(*Connection)

	network := &NetworkDescriptor{
		UUID: d.Id(),
	}

	if err := network.Load(c); err != nil {
		if xenErr, ok := err.(*xenAPI.Error); ok {
			if xenErr.Code() == xenAPI.ERR_UUID_INVALID {
				return false, nil
			}
		}

		return false, err
	}

	return true, nil
}