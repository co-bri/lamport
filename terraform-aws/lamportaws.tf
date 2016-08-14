
/*
*    Copyright (c) 2016 Brian McKean
*
*/

/*
*   Spins up three of the smallest server for use with lamport distributed
*   computing project
*   This script spins them up in AWS
*/

variable "num_nodes" {
    default = "3"
}


provider "aws" {
    access_key = "${var.access_key}"
    secret_key = "${var.secret_key}"
    region = "${var.region}"
}


resource "aws_security_group" "allow_all" {
  name = "allow_all"
  description = "Allow all inbound traffic"

  ingress {
      from_port = 0
      to_port = 0
      protocol = "-1"
      cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
      from_port = 0
      to_port = 0
      protocol = "-1"
      cidr_blocks = ["0.0.0.0/0"]
  }
  provisioner "local-exec" {
    command = "echo removing previous lamport_ips.txt"
    command = "rm -f lamport_ips.txt"
    command = "echo removing previous lamport_hosts.txt"
    command = "rm -f lamport_hosts.txt"
  }

}

resource "aws_instance" "lamport" {
  count = "${var.num_nodes}"
  ami           = "ami-fce3c696"
  instance_type = "t2.micro"
  tags {
        Name = "${format("LamportNode-%03d",count.index)}"
    }
  key_name = "${var.key_name}"
  provisioner "local-exec" {
    command = "echo ${self.public_dns} ${self.public_ip} >> lamport_hosts.txt"
    }
  provisioner "local-exec" {
    command = "echo ${self.tags.Name} ${self.public_dns} ${self.public_ip} >> lamport_ips.txt"
    }
  security_groups = ["${aws_security_group.allow_all.name}"]
}

