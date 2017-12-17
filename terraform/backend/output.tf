# Used for configuring ELBs.
//output "instance_ids" {
//  value = ["${aws_instance.instance.*.id}"]
//}

output "instance_internal_ip" {
  value = "${aws_instance.instance.private_ip}"
}

output "public_ips" {
  value = ["${aws_instance.instance.*.public_ip}"]
}
