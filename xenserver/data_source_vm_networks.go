package xenserver

import (
	"github.com/hashicorp/terraform/helper/schema"
	"log"
	"time"
)

const (
	vmNetworksVmUUID       = "vm_uuid"
	vmNetworksIp           = "ip"
	vmNetworksIpv6         = "ipv6"
	vmNetworksStartupDelay = "startup_delay"
)

func dataSourceVmNetworks() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceVmNetworksRead,
		Schema: map[string]*schema.Schema{
			vmNetworksVmUUID: &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			vmNetworksIp: &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeList,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},
			},
			vmNetworksIpv6: &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeList,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},
			},
			vmNetworksStartupDelay: &schema.Schema{
				Type:     schema.TypeFloat,
				Optional: true,
				Default:  0,
			},
		},
	}
}

func dataSourceVmNetworksRead(d *schema.ResourceData, meta interface{}) (err error) {
	c := meta.(*Connection)

	vm := &VMDescriptor{
		UUID: d.Get(vmNetworksVmUUID).(string),
	}
	if err = vm.Load(c); err != nil {
		return
	}

	if delay, ok := d.Get(vmNetworksStartupDelay).(float64); ok && delay > 0 {
		var vmmetrics VMMetrics
		if vmmetrics, err = vm.Metrics(c); err != nil {
			return
		}

		now := time.Now()
		diff := now.Sub(vmmetrics.StartTime).Seconds()

		if delay > diff {
			sleep := time.Duration(delay-diff) * time.Second
			time.Sleep(sleep)
		}
	}

	var metrics VMGuestMetrics

	if metrics, err = vm.GuestMetrics(c); err != nil {
		return
	}

	d.SetId(metrics.UUID)

	log.Printf("[DEBUG] Id is %s\n", d.Id())
	log.Println("[DEBUG] Networks: ", metrics.Networks)

	ipNetworks := make([][]string, 0)
	ipv6Networks := make([][]string, 0)

	for _, network := range metrics.Networks {
		ipNetworks = append(ipNetworks, network["ip"])
		ipv6Networks = append(ipv6Networks, network["ipv6"])
	}

	d.Set(vmNetworksIp, ipNetworks)
	d.Set(vmNetworksIpv6, ipv6Networks)

	return nil
}
