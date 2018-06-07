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

data "xenserver_vm_networks" "interfaces" {
 vm_uuid = "${vm_uuid}"
}

data "template_file" "interfaces_written" {
  template = "$${ips}"

  vars {
    ips = "${jsonencode(data.xenserver_vm_networks.interfaces.ip)}"
  }
}

resource "null_resource" "interfaces_file" {
  triggers {
    content = "${data.template_file.interfaces_written.rendered}"
  }

  provisioner "local-exec" {
    command = <<-EOC
      tee ${path.cwd}/output.json <<EOF
      ${data.template_file.interfaces_written.rendered}
      EOF
      EOC
  }
}

output "vm_main_ip" {
    value = "${${element(data.xenserver_vm_networks.interfaces.ip[0],0)}}"
}
