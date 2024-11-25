provider "aws" {
  alias               = "eu-west-1"
  region              = "eu-west-1"
  allowed_account_ids = [local.account_id]
  default_tags { tags = local.tags }
}

module "config_eu-west-1" {
  source    = "../modules/config_region"
  providers = { aws = aws.eu-west-1 }
  account   = local.account
}
