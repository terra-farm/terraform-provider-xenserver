package xenserver

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceXenServerSR() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceXenServerSRRead,

		Schema: map[string]*schema.Schema{
			"name_label": &schema.Schema{
				Type:        schema.TypeString,
				Description: "The human readable name of the storage repository",
				Required:    true,
			},
		},
	}
}

func dataSourceXenServerSRRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Connection)

	nameLabel, nameLabelOk := d.GetOk("name_label")

	if !nameLabelOk {
		return fmt.Errorf("name_label must be provided")
	}

	if srs, err := c.client.SR.GetByNameLabel(c.session, nameLabel.(string)); err == nil {
		found := false
		for _, sr := range srs {
			record, err := c.client.SR.GetRecord(c.session, sr)
			if err != nil {
				break
			}

			d.SetId(record.UUID)

			found = true
			break
		}

		if !found {
			return fmt.Errorf("Matching SR not found")
		}
	}

	return nil
}
