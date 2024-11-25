provider "aws" {
  alias               = "eu-central-1"
  region              = "eu-central-1"
  allowed_account_ids = [local.account_id]
  default_tags { tags = local.tags }
}

module "config_eu-central-1" {
  source    = "../modules/config_region"
  providers = { aws = aws.eu-central-1 }
  account   = local.account
}
