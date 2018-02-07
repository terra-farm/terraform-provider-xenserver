---
title: "xenserver_vm_data_metrics"
---

Provides metrics reported by guest tools.

## Parameters
* [in] vm_uuid - UUID of the VM to obtain metrics from
* [out] disks
* [out] networks
* [out] memory
* [out] os_version
* [out] pv_driver_version
* [out] is_pv_driver_present
* [out] can_use_hotplug_vbd
* [out] can_use_hotplug_vif
* [out] is_live
* [out] last_updated

## Example Usage

```hcl
data "xenserver_vm_guest_metrics" "interfaces" {
 vm_uuid = "${vm_uuid}"
}

output "vm_main_interface" {
    value = "${data.xenserver_vm_guest_metrics.interfaces.networks[0]}"
}
```
