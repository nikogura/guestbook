# Used for configuring ELBs.
output "instance_ids" {
  value = ["${aws_instance.instance.*.id}"]
}

output "public_ips" {
  value = ["${aws_instance.instance.*.public_ip}"]
}
