# To execute with `credentials.tfvars` in the `examples` folder:
# terraform init && terraform plan --var-file=../credentials.tfvars && terraform apply --var-file=../credentials.tfvars && terraform destroy --var-file=../credentials.tfvars -force

variable "url" {}
variable "username" {}
variable "password" {}

provider "xenserver" {
  url      = "${var.url}"
  username = "${var.username}"
  password = "${var.password}"
}

resource "xenserver_vm" "ansiblevm" {
  base_template_name = "CentOS 7.4 Template"
  name_label = "Ansible VM"
  static_mem_min = 2147483648
  static_mem_max = 2147483648
  dynamic_mem_min = 2147483648
  dynamic_mem_max = 2147483648
  vcpus = 1
  boot_order = "c"
  hard_drive {
    is_from_template = true
    user_device = "0"
  }
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
  network_interface {
    network_uuid = "9318f024-7937-870f-32c7-171040f1fbd8"
    device = 1
    mtu = 1500
    mac = ""
  }
}

data "xenserver_vm_networks" "ansible_interfaces" {
 vm_uuid = "${xenserver_vm.ansiblevm.id}"
 startup_delay = 10 # wait for VM to boot for 10 seconds
}

resource "null_resource" "provision_ansiblevm" {
  provisioner "local-exec" {
    command = 
"sleep 120; ANSIBLE_HOST_KEY_CHECKING=False ansible-playbook -i '${element(data.xenserver_vm_networks.ansible_interfaces.ip[0],0)},' --extra-vars 'provisioned_host_name=ansible.local' ${path.cwd}/provision.yaml"
  }
}
