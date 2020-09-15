package xenserver

import (
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceXenServerPifs() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceXenServerPifsRead,
		Schema: map[string]*schema.Schema{
			"uuids": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceXenServerPifsRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Connection)

	pifUUIDs := make([]string, 0)

	if pifs, err := c.client.PIF.GetAllRecords(c.session); err == nil {
		for _, pif := range pifs {
			pifUUIDs = append(pifUUIDs, pif.UUID)
		}
	}

	d.SetId(time.Now().UTC().String())
	d.Set("uuids", pifUUIDs)

	return nil
}
