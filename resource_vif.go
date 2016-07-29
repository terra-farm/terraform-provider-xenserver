package main

import (
	"github.com/amfranz/go-xen-api-client"
	"github.com/hashicorp/terraform/helper/schema"
)

const (
	vifSchemaNetworkNameLabel        = "network_name_label"
)

func resourceVIF() *schema.Resource {
	return &schema.Resource{
		Create: resourceVIFCreate,
		Read:   resourceVIFRead,
		Update: resourceVIFUpdate,
		Delete: resourceVIFDelete,
		Exists: resourceVIFExists,

		Schema: map[string]*schema.Schema{
			vifSchemaNetworkNameLabel: &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceVIFCreate(d *schema.ResourceData, m interface{}) error {
	c := m.(*Connection)

	networks, err := c.client.Network.GetByNameLabel(c.session, d.Get(vifSchemaNetworkNameLabel).(string))
	if err != nil {
		return err
	}

	network:= networks[0]

	vif := xenAPI.VIFRecord{
		Network: network,
	}

	vifRef, err := c.client.VIF.Create(c.session, vif)
	if err != nil {
		return err
	}

	vifObject, err := c.client.VIF.GetRecord(c.session, vifRef)
	if err != nil {
		return err
	}

	d.SetId(vifObject.UUID)

	return nil
}

func resourceVIFRead(d *schema.ResourceData, m interface{}) error {
	c := m.(*Connection)

	vifRef, err := c.client.VIF.GetByUUID(c.session, d.Id())
	if err != nil {
		if xenErr, ok := err.(*xenAPI.Error); ok {
			if xenErr.Code() == xenAPI.ERR_UUID_INVALID {
				d.SetId("")
				return nil
			}
		}

		return err
	}

	vif, err := c.client.VIF.GetRecord(c.session, vifRef)
	if err != nil {
		return err
	}

	net, err := c.client.Network.GetRecord(c.session, vif.Network)
	if err != nil {
		return err
	}

	d.Set(vifSchemaNetworkNameLabel, net.NameLabel)

	return nil
}

func resourceVIFUpdate(d *schema.ResourceData, m interface{}) error {
	c := m.(*Connection)

	vifRef, err := c.client.VIF.GetByUUID(c.session, d.Id())
	if err != nil {
		if xenErr, ok := err.(*xenAPI.Error); ok {
			if xenErr.Code() == xenAPI.ERR_UUID_INVALID {
				d.SetId("")
				return nil
			}
		}

		return err
	}

	vif, err := c.client.VIF.GetRecord(c.session, vifRef)
	if err != nil {
		return err
	}

	d.Partial(true)

	netLabel := d.Get(vifSchemaNetworkNameLabel).(string)

	net, err := c.client.Network.GetRecord(c.session, vif.Network)
	if err != nil {
		return err
	}

	if netLabel != net.NameLabel {

		networks, err := c.client.Network.GetByNameLabel(c.session, netLabel)
		if err != nil {
			return err
		}

		network := networks[0]

		err = c.client.VIF.Destroy(c.session, vifRef)
		if err != nil {
			return err
		}

		vif.Network = network

		vifRef, err = c.client.VIF.Create(c.session, vif)
		if err != nil {
			return err
		}

		d.SetPartial(vifSchemaNetworkNameLabel)

	}

	return nil
}

func resourceVIFDelete(d *schema.ResourceData, m interface{}) error {
	c := m.(*Connection)

	vifRef, err := c.client.VIF.GetByUUID(c.session, d.Id())
	if err != nil {
		if xenErr, ok := err.(*xenAPI.Error); ok {
			if xenErr.Code() == xenAPI.ERR_UUID_INVALID {
				d.SetId("")
				return nil
			}
		}

		return err
	}

	err = c.client.VIF.Destroy(c.session, vifRef)
	if err != nil {
		return err
	}

	return nil
}

func resourceVIFExists(d *schema.ResourceData, m interface{}) (bool, error)  {
	c := m.(*Connection)

	_, err := c.client.VIF.GetByUUID(c.session, d.Id())
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