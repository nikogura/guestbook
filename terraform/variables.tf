variable "public_key_path" {
  description = "Enter the path to the SSH Public Key to add to AWS."
  default = "~/.ssh/id_rsa.pub"
}

variable "aws_region" {
  description = "AWS region to launch servers."
  default     = "us-east-1"
}

variable "home_ip" {
  description = "My home IP"
  default = "73.158.93.80/32"
}












