# To execute with `credentials.tfvars` in the `examples` folder:
# terraform init && terraform plan --var-file=../credentials.tfvars && terraform apply --var-file=../credentials.tfvars && terraform destroy --var-file=../credentials.tfvars -force

variable "url" {}
variable "username" {}
variable "password" {}
variable "vm_uuid" {}

provider "xenserver" {
  url      = "${var.url}"
  username = "${var.username}"
  password = "${var.password}"
}

resource "xenserver_vm" "newvm" {
  base_template_name = "CentOS 7.4 Template"
  name_label = "New VM"
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

data "xenserver_vm_networks" "new_interfaces" {
 vm_uuid = "${xenserver_vm.newvm.id}"
 startup_delay = 10 # wait for VM to boot for 10 seconds
}

data "template_file" "new_interfaces_written" {
  template = "$${ips}"

  vars {
    ips = "${jsonencode(data.xenserver_vm_networks.new_interfaces.ip)}"
  }
}

resource "null_resource" "new_interfaces_file" {
  triggers {
    content = "${data.template_file.new_interfaces_written.rendered}"
  }

  provisioner "local-exec" {
    command = <<-EOC
      tee ${path.cwd}/newvm_output.json <<EOF
      ${data.template_file.new_interfaces_written.rendered}
      EOF
      EOC
  }
}

output "new_vm_main_ip" {
    value = "${${element(data.xenserver_vm_networks.new_interfaces.ip[0],0)}}"
}
