---
title: "xenserver_pif"
---

Provides information about a physical network interface (PIF) of a XenServer host specified by the interface name or whether it is the management interface.

## Example Usage

```hcl
data "xenserver_pif" "eth0" {
  device = "eth0"
}

data "xenserver_pif" "management" {
  management = true
}

resource "xenserver_vm" "demo-vm" {
  // ...
  network_interface {
    network_uuid = "${data.xenserver_pif.management.network_uuid}"
    device = 0
  }
  network_interface {
    network_uuid = "${data.xenserver_pif.eth0.network_uuid}"
    device = 1
  }
  // ...
}
```
