provider "aws" {
  region = "${var.aws_region}"
}

resource "aws_vpc" "vpc_main" {
  cidr_block           = "10.0.0.0/16"
  enable_dns_support   = true
  enable_dns_hostnames = true

  tags {
    Name = "Main VPC"
  }
}

resource "aws_key_pair" "auth" {
  key_name   = "default"
  public_key = "${file(var.public_key_path)}"
}

resource "aws_internet_gateway" "default" {
  vpc_id                = "${aws_vpc.vpc_main.id}"
}

resource "aws_route" "internet_access" {
  route_table_id          = "${aws_vpc.vpc_main.main_route_table_id}"
  destination_cidr_block  = "0.0.0.0/0"
  gateway_id              = "${aws_internet_gateway.default.id}"
}

# Create a public subnet to launch our load balancers
resource "aws_subnet" "public" {
  vpc_id                  = "${aws_vpc.vpc_main.id}"
  cidr_block              = "10.0.1.0/24" # 10.0.1.0 - 10.0.1.255 (256)
  map_public_ip_on_launch = true
  #availability_zone       = "${var.aws_availability_zone}"
}

# Create a private subnet for web servers
resource "aws_subnet" "frontend" {
  vpc_id                  = "${aws_vpc.vpc_main.id}"
  cidr_block              = "10.0.2.0/24"
  map_public_ip_on_launch = true
  #availability_zone       = "${var.aws_availability_zone}"
}

# Create a private subnet for app servers
resource "aws_subnet" "app" {
  vpc_id                  = "${aws_vpc.vpc_main.id}"
  cidr_block              = "10.0.3.0/24"
  map_public_ip_on_launch = true
  #availability_zone       = "${var.aws_availability_zone}"
}

# Create a private subnet for backend
resource "aws_subnet" "backend" {
  vpc_id                  = "${aws_vpc.vpc_main.id}"
  cidr_block              = "10.0.4.0/24"
  map_public_ip_on_launch = true
  #availability_zone       = "${var.aws_availability_zone}"
}

# A security group for the ELB so it is accessible via the web
resource "aws_security_group" "elb" {
  name        = "sec_group_elb"
  description = "Security group for public facing ELBs"
  vpc_id      = "${aws_vpc.vpc_main.id}"

  # HTTP access from anywhere
  ingress {
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  # 8080 access from anywhere
  ingress {
    from_port   = 8080
    to_port     = 8080
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  # Outbound internet access
  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

# security group to access the frontend from the load balancers
resource "aws_security_group" "frontend" {
  name        = "sec_group_frontend"
  description = "Security group for frontend servers"
  vpc_id      = "${aws_vpc.vpc_main.id}"

  # SSH access from anywhere
  ingress {
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  # HTTP access from the VPC
  ingress {
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = ["${aws_subnet.public.cidr_block}"]
  }

  # Outbound internet access
  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

//# security group to access the app from the frontend and elsewhere
resource "aws_security_group" "app" {
  name        = "sec_group_app"
  description = "Security group for app servers"
  vpc_id      = "${aws_vpc.vpc_main.id}"

  # SSH access from anywhere
  ingress {
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  # HTTP access from the load balancer
  ingress {
    from_port   = 8080
    to_port     = 8080
    protocol    = "tcp"
    cidr_blocks = ["${aws_subnet.app.cidr_block}"]
  }

  # Outbound internet access
  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

# security group to access the backend
resource "aws_security_group" "backend" {
  name        = "sec_group_backend"
  description = "Security group for backend database servers"
  vpc_id      = "${aws_vpc.vpc_main.id}"

  # SSH access else we can't provision.  Yup.  It's a backdoor.  Be different if we had a VPN, or if I was using a fancier provisioner.
  ingress {
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    #cidr_blocks = ["${aws_subnet.app.cidr_block}"]
    cidr_blocks = ["${var.home_ip}", "${aws_subnet.frontend.cidr_block}"]
  }

  # Access from the app servers
  ingress {
    from_port   = 5432
    to_port     = 5432
    protocol    = "tcp"
    cidr_blocks = ["${aws_subnet.app.cidr_block}"]
  }

  # Outbound internet access
  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}


module "frontend" {
  source                 = "./frontend"
  subnet_id              = "${aws_subnet.frontend.id}"
  key_pair_id            = "${aws_key_pair.auth.id}"
  security_group_id      = "${aws_security_group.frontend.id}"

  count                  = 2
  group_name             = "frontend"
  downstream_server      = "${aws_elb.app.dns_name}"
}

module "app" {
  source                 = "./app"
  subnet_id              = "${aws_subnet.app.id}"
  key_pair_id            = "${aws_key_pair.auth.id}"
  security_group_id      = "${aws_security_group.app.id}"

  count                  = 2
  group_name             = "app"
  downstream_server      = "${module.backend.instance_internal_ip}"
}

module "backend" {
  source                 = "./backend"
  subnet_id              = "${aws_subnet.backend.id}"
  key_pair_id            = "${aws_key_pair.auth.id}"
  security_group_id      = "${aws_security_group.backend.id}"

  count                  = 1
  disk_size              = 10
  group_name             = "backend"
  instance_type          = "t2.micro"

  app_network        = "${aws_subnet.app.cidr_block}"
}

data "aws_availability_zones" "available" {}

# Frontend ELB
resource "aws_elb" "frontend" {
  name = "elb-frontend"

  subnets         = ["${aws_subnet.public.id}", "${aws_subnet.frontend.id}"]
  security_groups = ["${aws_security_group.elb.id}"]
  instances       = ["${module.frontend.instance_ids}"]

  listener {
    instance_port     = 80
    instance_protocol = "http"
    lb_port           = 80
    lb_protocol       = "http"
  }

  health_check {
    healthy_threshold   = 2
    unhealthy_threshold = 2
    timeout             = 3
    target              = "HTTP:80/"
    interval            = 30
  }
}

# App ELB
resource "aws_elb" "app" {
  name = "elb-app"

  subnets         = ["${aws_subnet.app.id}"]
  security_groups = ["${aws_security_group.elb.id}"]
  instances       = ["${module.app.instance_ids}"]

  listener {
    instance_port     = 8080
    instance_protocol = "http"
    lb_port           = 8080
    lb_protocol       = "http"
  }

  health_check {
    healthy_threshold   = 2
    unhealthy_threshold = 2
    timeout             = 3
    target              = "HTTP:8080/guestbook/healthcheck"
    interval            = 30
  }
}

# Not going to mess with a load balanced database at this point.
# Certainly it can be done, but how really depends on business needs.  Can we use the RDS?  Is cost a factor?  DynamoDB would be fine in this particular case, but MySQL or Postgres was the requirement.

# Private ELB for backend
//resource "aws_elb" "backend" {
//  name = "elb-backend"
//
//  subnets         = ["${aws_subnet.backend.id}"]
//  security_groups = ["${aws_security_group.backend.id}"]
//  instances       = ["${module.backend.instance_ids}"]
//  internal        = true
//
//  listener {
//    instance_port     = 5432
//    instance_protocol = "tcp"
//    lb_port           = 5432
//    lb_protocol       = "tcp"
//  }
//
//  health_check {
//    healthy_threshold   = 1
//    unhealthy_threshold = 0
//    timeout             = 3
//    target              = "TCP:5432"
//    interval            = 30
//  }
//}