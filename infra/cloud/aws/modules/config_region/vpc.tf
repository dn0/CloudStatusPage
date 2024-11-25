module "vpc" {
  source  = "terraform-aws-modules/vpc/aws"
  version = "~> 5.14.0"

  name           = "vpc"
  cidr           = local.region_cfg.cidr
  azs            = local.region_cfg.azs
  public_subnets = local.region_cfg.public_subnets

  create_igw                      = true
  create_egress_only_igw          = true
  enable_nat_gateway              = false
  create_elasticache_subnet_group = false
  create_redshift_subnet_group    = false

  default_security_group_ingress = [
    {
      from_port   = 8
      to_port     = 0
      protocol    = "icmp"
      cidr_blocks = "10.0.0.0/8,172.16.0.0/12,192.168.0.0/16"
      description = "ICMP ping"
    }
  ]

  default_security_group_egress = [
    {
      protocol    = "-1"
      from_port   = 0
      to_port     = 0
      cidr_blocks = "0.0.0.0/0"
      description = "Outbound traffic"
    }
  ]

  tags = {
    cost-center = "mon-agent"
  }
}
