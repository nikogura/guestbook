variable "private_key_path" {
  description = "Enter the path to the SSH Private Key to run provisioner."
  default = "~/.ssh/id_rsa"
}

variable "aws_ami" {
  default = "ami-c29e1cb8"
}

variable "disk_size" {
  default = 8
}

variable "count" {
  default = 1
}

variable "group_name" {
  description = "Group name becomes the base of the hostname of the instance"
}

variable "aws_region" {
  description = "AWS region to launch servers."
  default     = "us-west-1"
}

variable "instance_type" {
  description = "Default Instance Type"
  default     = "t2.micro"
}

variable "subnet_id" {
  description = "ID of the AWS VPC subnet to use"
}

variable "key_pair_id" {
  description = "ID of the keypair to use for SSH"
}

variable "security_group_id" {
  description = "ID of the VPC security group to use for network"
}

variable "downstream_server" {
  description = "DNS name of downstream server or ELB"
}
