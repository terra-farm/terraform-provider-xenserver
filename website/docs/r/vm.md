---
title: "xenserver_vm"
---

Provides a XenServer virtual machine resource. This can be used to create, modify, and delete virtual machines.

## Example Usage

```hcl
resource "xenserver_vm" "web" {
    name_label = "web"
    base_template_name = "<desired template>"
    static_mem_min = 2147483648 # 2GB
    static_mem_max = 2147483648
    dynamic_mem_min = 2147483648
    dynamic_mem_max = 2147483648
    boot_order = "cdn"
    network_interface {
        network_uuid = "<uuid>"
        mac = "<desired-mac>"
        mtu = 1500
        device = 0
    }
    network_interface {
        network_uuid = "<desired network>"
        mtu = 1500
        device = 1
    }
    vcpus = 2
    cdrom {
      vdi_uuid = "<iso uuid>"
    }
    hard_drive {
      vdi_uuid = "<desired vdi>"
    }
    other_config {
        auto_poweron = "true"
    }
}
```

## Argument Reference

The following arguments are supported:

* `name_label` - (Required) The name given for this VM.
* `base_template_name` - 
* `static_mem_min` - 
* `static_mem_max` - 
* `dynamic_mem_min` - 
* `boot_order` - 
* `vcpus` - 

The `network_interface` block supports:

* `network_uuid` -
* `mtu` -
* `device` -

The `cdrom` block supports:

* `vdi_uuid` - 

The `hard_drive` block supports:

* `vdi_uuid` - 

The `other_config` block sets any number of given key-value pairs in the VM's `other-config` map.

## Attributes Reference

The following attributes are exported:

* `id` - The instance ID.
