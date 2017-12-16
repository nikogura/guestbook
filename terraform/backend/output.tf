# Used for configuring ELBs.
//output "instance_ids" {
//  value = ["${aws_instance.instance.*.id}"]
//}

output "instance_internal_ip" {
  value = "${aws_instance.instance.private_ip}"
}

