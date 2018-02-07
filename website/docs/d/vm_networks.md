---
title: "xenserver_vm_networks"
---

Query information on IP addresses reported by guest tools.

## Parameters
* [in] vm_uuid - UUID of the VM to query information about
* [in] [optional] startup_delay - how many seconds should be passed after VM start before information could be queried. Useful when new VM was just created to wait until it boots and guest tools are started
* [out] ip array of arrays with ip addresses where first index corresponds to interface and second to address
* [out] ipv6 array of arrays with IPv6 addresses where first index corresponds to interface and second to address

## Example Usage

```hcl
data "xenserver_vm_networks" "interfaces" {
 vm_uuid = "${vm_uuid}"
}

output "vm_main_ip" {
    value = "${${element(data.xenserver_vm_networks.interfaces.ip[0],0)}}"
}
```
