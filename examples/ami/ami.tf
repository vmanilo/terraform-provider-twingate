# define or use an existing VPC
module "demo_vpc" {
  source  = "terraform-aws-modules/vpc/aws"
  version = "~> 5.8"

  name = "demo_vpc"
  cidr = "10.0.0.0/16"

  azs                  = ["us-east-1a"]
  private_subnets      = ["10.0.1.0/24"]
  public_subnets       = ["10.0.2.0/24"]
  enable_dns_hostnames = true
  enable_nat_gateway   = true
}

# If you use an existing Security group, the Connector requires egress traffic enabled but does not require ingress
module "demo_sg" {
  source             = "terraform-aws-modules/security-group/aws"
  version            = "~> 5.1"
  vpc_id             = module.demo_vpc.vpc_id
  name               = "demo_security_group"
  egress_cidr_blocks = ["0.0.0.0/0"]
  egress_rules       = ["all-tcp", "all-udp", "all-icmp"]
}

# spin off a ec2 instance from Twingate AMI and configure tokens in user_data
module "ec2_tenant_connector" {
  source  = "terraform-aws-modules/ec2-instance/aws"
  version = "~> 5.6"

  name                   = "demo_connector"
  user_data              = <<-EOT
    #!/bin/bash
    set -e
    mkdir -p /etc/twingate/
    {
      echo TWINGATE_URL="https://${var.network}.twingate.com"
      echo TWINGATE_ACCESS_TOKEN="${twingate_connector_tokens.aws_connector_tokens.access_token}"
      echo TWINGATE_REFRESH_TOKEN="${twingate_connector_tokens.aws_connector_tokens.refresh_token}"
    } > /etc/twingate/connector.conf
    sudo systemctl enable --now twingate-connector
  EOT
  ami                    = data.aws_ami.latest.id
  instance_type          = "t3a.micro"
  vpc_security_group_ids = [module.demo_sg.security_group_id]
  subnet_id              = module.demo_vpc.private_subnets[0]

  root_block_device = [
    {
      encrypted = true
    }
  ]
}
