# Terraform provider for XenServer
[![GitHub release](http://img.shields.io/github/release/mikljohansson/terraform-provider-xenserver.svg)](https://github.com/mikljohansson/terraform-provider-xenserver/releases)
[![Travis CI](https://img.shields.io/travis/mikljohansson/terraform-provider-xenserver/master.svg)](https://travis-ci.org/mikljohansson/terraform-provider-xenserver)

## Usage

### Provider Configuration

Download and place the `terraform-provider-xenserver` binary into your $PATH

```
provider "xenserver" {
    url = "${var.xenserver_url}"
    username = "${var.xenserver_username}"
    password = "${var.xenserver_password}"
}
```

The following arguments are supported.

* `url` - (Required) The URL to the XenAPI endpoint, typically "https://<XenServer Management IP>"
* `username` - (Required) The username to use to authenticate to XenServer.
* `password` - (Required) The password to use to authenticate to XenServer.

### Resource Configuration

#### `xenserver_vm`

```
resource "xenserver_vm" "myvm" {
    name_label = "My VM"
    base_template_name = "centos-7-large"
    xenstore_data {
        hostname = "myvm.example.com"
        ip = "192.168.1.20"
    }
}
```

The following arguments are supported.

* `name_label` - (Required) Name of VM.
* `base_template_name` - (Required) Name VM template to instantiate.
* `xenstore_data` - (Optional) Configuration made available inside the VM as "vm-data/key=value" using the xenstore-read utility.
