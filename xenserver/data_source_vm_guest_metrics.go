package xenserver

import (
	"github.com/hashicorp/terraform/helper/schema"
	"log"
)

const (
	vmGuestMetricsVmUUID            = "vm_uuid"
	vmGuestMetricsDisks             = "disks"
	vmGuestMetricsNetworks          = "networks"
	vmGuestMetricsMemory            = "memory"
	vmGuestMetricsOSVersion         = "os_version"
	vmGuestMetricsPVDriversVersion  = "pv_driver_version"
	vmGuestMetricsPVDriversDetected = "is_pv_driver_present"
	vmGuestMetricsCanUseHotPlugVbd  = "can_use_hotplug_vbd"
	vmGuestMetricsCanUseHotPlugVif  = "can_use_hotplug_vif"
	vmGuestMetricsLive              = "is_live"
	vmGuestMetricsLastUpdated       = "last_updated"
)

func dataSourceVmGuestMetrics() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceVmGuestMetricsRead,
		Schema: map[string]*schema.Schema{
			vmGuestMetricsVmUUID: &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			vmGuestMetricsDisks: &schema.Schema{
				Type:     schema.TypeMap,
				Computed: true,
			},
			vmGuestMetricsNetworks: &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip": &schema.Schema{
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"ipv6": &schema.Schema{
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			vmGuestMetricsMemory: &schema.Schema{
				Type:     schema.TypeMap,
				Computed: true,
			},
			vmGuestMetricsOSVersion: &schema.Schema{
				Type:     schema.TypeMap,
				Computed: true,
			},
			vmGuestMetricsPVDriversVersion: &schema.Schema{
				Type:     schema.TypeMap,
				Computed: true,
			},
			vmGuestMetricsPVDriversDetected: &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
			vmGuestMetricsCanUseHotPlugVbd: &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			vmGuestMetricsCanUseHotPlugVif: &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			vmGuestMetricsLive: &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
			vmGuestMetricsLastUpdated: &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceVmGuestMetricsRead(d *schema.ResourceData, meta interface{}) (err error) {
	c := meta.(*Connection)

	vm := &VMDescriptor{
		UUID: d.Get(vmGuestMetricsVmUUID).(string),
	}
	if err = vm.Load(c); err != nil {
		return
	}

	var metrics VMGuestMetrics

	if metrics, err = vm.GuestMetrics(c); err != nil {
		return
	}

	d.SetId(metrics.UUID)

	log.Printf("[DEBUG] Id is %s\n", d.Id())
	log.Println("[DEBUG] Networks: ", metrics.Networks)

	d.Set(vmGuestMetricsDisks, metrics.Disks)
	d.Set(vmGuestMetricsNetworks, metrics.Networks)
	d.Set(vmGuestMetricsMemory, metrics.Memory)
	d.Set(vmGuestMetricsOSVersion, metrics.OSVersion)
	d.Set(vmGuestMetricsPVDriversVersion, metrics.PVDriversVersion)
	d.Set(vmGuestMetricsPVDriversDetected, metrics.PVDriversDetected)
	d.Set(vmGuestMetricsCanUseHotPlugVbd, metrics.CanUseHotplugVbd)
	d.Set(vmGuestMetricsCanUseHotPlugVif, metrics.CanUseHotplugVif)
	d.Set(vmGuestMetricsLive, metrics.Live)
	d.Set(vmGuestMetricsLastUpdated, metrics.LastUpdated.String())

	return nil
}

/*
func networksToSchemaList(networks []map[string][]string) []interface{} {



}*/
