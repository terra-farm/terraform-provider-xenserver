---
title: XenServer Provider
---

The XenServer provider is used to interact with the resources supported by XenServer.
The provider needs to be configured with the proper credentials before it can be used.

## Installation

In your home folder, create a file `~/.terraformrc` with these contents:

```hcl
providers {
  xenserver = "<full path>/terraform-provider-xenserver"
}
```

NOTE: this is not yet updated for Terraform 0.11.

## Example Usage

```hcl
# Configure the XenServer Provider
provider "xenserver" {
  url      = "<xen-server-url>"
  username = "<user-name>"
  password = "<password>"
}

// Create a new virtual machine
resource "xenserver_vm" "test" {
  base_template_name = "CentOS 7.4 Template"
  name_label = "test"
  static_mem_min = 8589934592
  static_mem_max = 8589934592
  dynamic_mem_min = 8589934592
  dynamic_mem_max = 8589934592
  vcpus = 1
  boot_order = "c"

  hard_drive {
    is_from_template = true
    user_device = "0"
  } # Template VM HDD
  cdrom {
    is_from_template = true
    user_device = "3"
  }
  network_interface {
    network_uuid = "92467b56-21a7-dfdd-b412-978181a69f32"
    device = 0
    mtu = 1500
    mac = ""
    other_config {
        ethtool-gso = "off"
        ethtool-ufo = "off"
        ethtool-tso = "off"
        ethtool-sg = "off"
        ethtool-tx = "off"
        ethtool-rx = "off"
    }
  }
}
```

## XenServer versions

Both backward and forward compatibility with the XenApi is mostly defined by the
[go-xen-api-client](https://github.com/amfranz/go-xen-api-client) Go library.

Tested succesfully against:
* XenServer 7.2

## Authentication

Authentication against the XenApi happens with username and password combination.
To protect your credentials for going in clear text over the wire, it is advised
to connect to SSL/TLS endpoints.

```hcl
provider "kubernetes" {
  url      = "https://104.196.242.174"
  username = "XenMaster"
  password = "XenInfraAsCode"
}
```

## Argument Reference

The following arguments are supported:

* `url` - (Required) the XenApi endpoint of your XenServer or XenServer pool.
* `username` - (Required) The username to use for HTTP basic authentication when accessing
  the XenApi endpoint.
* `password` - (Required) The password to use for HTTP basic authentication when accessing
  the XenApi endpoint.
