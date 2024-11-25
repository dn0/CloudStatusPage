terraform {
  required_version = "~> 1.8.0"

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "5.67.0"
    }
  }
}

locals {
  account    = "cloudstatus-probe-p"
  account_id = "123456"

  tags = {
    terraform   = "true"
    environment = "production"
  }
}

module "budget" {
  source       = "../modules/budget"
  providers    = { aws = aws.eu-central-1 }
  limit_amount = 330
}

module "config_global" {
  source     = "../modules/config_global"
  providers  = { aws = aws.eu-central-1 }
  account    = local.account
  account_id = local.account_id
}
