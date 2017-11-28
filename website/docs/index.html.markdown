---
layout: "xenserver"
page_title: "Provider: XenServer"
sidebar_current: "docs-xenserver-index"
description: |-
  The XenServer provider is used to interact with the resources supported by XenServer. 
  The provider needs to be configured with the proper credentials before it can be used.
---

# XenServer Provider

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
resource "xenserver_vm" "my_machine" {
    # ...
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

## Reference

### Data Sources

* [xenserver_pifs](d/pifs.html.markdown)

### Resources

* [xenserver_network](r/network.html.markdown)
* [xenserver_sr](r/sr.html.markdown)
* [xenserver_vbd](r/vbd.html.markdown)
* [xenserver_vdi](r/vdi.html.markdown)
* [xenserver_vif](r/vif.html.markdown)
* [xenserver_vlan](r/vlan.html.markdown)
* [xenserver_vm](r/vm.html.markdown)
