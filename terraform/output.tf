# Public Load Balancers

output "frontend_url" {
  value = "http://${aws_elb.frontend.dns_name}/guestbook"

}

output "frontend_server_ips" {
  value = "${module.frontend.public_ips}"
}

output "app_server_ips" {
  value = "${module.app.public_ips}"
}

output "backend_server_ips" {
  value = "${module.backend.public_ips}"
}
