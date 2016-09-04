# Terraform provider for XenServer

## Dependencies

    go get 'github.com/hashicorp/terraform'
    go get 'github.com/mborodin/go-xen-api-client'
    go get 'github.com/amfranz/go-xmlrpc-client'

### Usage example

#### Provider configuration
```
provider "xenserver" {
  url = "<xen-server-url>"
  username = "<user-name>"
  password = "<password>"
}
```
Arguments:
 * url - URL to XenServer
 * username - XenServer user allowed to execute APIs
 * password - user password

#### VM Creation
```
resource "xenserver_vdi" "vdi" {
  sr_uuid = "<sr uuid>"
  name_label = "test vdi"
  size = 1073741824 # 1GB
}

resource "xenserver_network" "net" {
  name_label = "vm only network"
  bridge = "xapi1"
  description = "this will create VM-only network"
}

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
        network_uuid = "${xenserver_network.net.id}"
        mtu = 1500
        device = 1
    }
    vcpus = 2
    cdrom {
      vdi_uuid = "<iso uuid>"
    }
    hard_drive {
      vdi_uuid = "${xenserver_vdi.vdi.id}"
    }

}
```
Arguments:
  * name_label - VM name
  * base_template_name - VM template
  * vcpus - Number of CPU's
  * static_mem_min - Minimal static memory (in bytes)
  * static_mem_max - Maximal static memory (in bytes)
  * dynamic_mem_min - Minimal dynamic memory (in bytes)
  * dynamic_mem_max - Maximal dynamic memory (in bytes)
  * boot_order (optional) - boot order. Use c for first bootable hard drive, d for CD-ROM, n for netboot. Example: cdn - boot from HD, CD, Network
  * network_interface (optional) - [Multiple possible] definition of network interface
  * hard_drive (optional) - [Multiple possible] connected hard drive
  * cdrom (optional) - [Multiple possible] connected cdrom
  * boot_parameters (optional) - TBD
  * installation_media_type (optional) - TDB
  * installation_media_location (optional) - TBD
  * cores_per_socket (optional) - CPU topology (default is 1 core per socket)
  * xenstore_data (optional) - Extra VM configuration data for in-VM use

Network interface schema:
  * network_uuid - Network UUID. Required either network name of network uuid
  * mac - Desired mac address
  * mtu - MTU
  * device - interface order
  * other_config - other configuration parameters map

Block device schema:
  * vdi_uuid - UUID of connected VDI
  * bootable (optional) - is device bootable. Default: *false*
  * mode (optional) - RW or RO. Default: *RW*

VDI Resource Schema
  * sr_uuid - SR that will hold newly created VDI image
  * name_label - VDI name
  * size - size (in bytes)
  * shared (optional) - should this VDI be shared among multiple VMs, Default: *false*
  * read_only (optional) - is this VDI is read-only. Default: *false*

Network Resource Schema:
  * name_label - Network Name
  * bridge - Bridge interface name
  * description - Network description
  * mtu - MTU
