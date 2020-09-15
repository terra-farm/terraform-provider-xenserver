package xenserver

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceXenServerPif() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceXenServerPifRead,

		Schema: map[string]*schema.Schema{
			"device": &schema.Schema{
				Type:        schema.TypeString,
				Description: "The machine-readable name of the physical interface (PIF) (e.g. eth0)",
				Optional:    true,
			},
			"management": &schema.Schema{
				Type:        schema.TypeBool,
				Description: "Indicates whether the control software is listening for connections on this physical interface",
				Optional:    true,
			},
			// Computed values
			"network": &schema.Schema{
				Type:        schema.TypeString,
				Description: "UUID of the virtual network to which this PIF is connected",
				Computed:    true,
			},
		},
	}
}

func dataSourceXenServerPifRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Connection)

	device, deviceOk := d.GetOk("device")
	management, managementOk := d.GetOk("management")

	if !deviceOk && !managementOk {
		return fmt.Errorf("One of device or management must be assigned")
	}

	if pifs, err := c.client.PIF.GetAllRecords(c.session); err == nil {
		found := false
		for _, pif := range pifs {
			if (!deviceOk || pif.Device == device) && (!managementOk || pif.Management == management) {
				d.SetId(pif.UUID)

				network := &NetworkDescriptor{
					NetworkRef: pif.Network,
				}
				if err = network.Query(c); err != nil {
					return err
				}
				d.Set("network", network.UUID)

				found = true
				break
			}
		}

		if !found {
			return fmt.Errorf("Matching PIF not found")
		}
	}

	return nil
}
