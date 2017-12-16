# Public Load Balancers

//output "backend_address" {
//  value = "${aws_elb.backend.dns_name}"
//}

output "frontend_url" {
  value = "http://${aws_elb.frontend.dns_name}/guestbook"

}

# Private Load Balancers

//output "db_loadbalancer" {
//  value = "${aws_elb.db_postgres.dns_name}"
//}