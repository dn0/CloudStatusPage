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
  account    = "cloudstatus-org"
  account_id = "123456"
  region     = "eu-central-1"

  tags = {
    terraform   = "true"
    environment = "org"
  }

  org_accounts = {
    mon-org = local.account_id
  }

  probe_accounts = {
    mon-probe-t = "123456"
    mon-probe-s = "234567"
    mon-probe-p = "345678"
  }

  all_accounts = merge(local.org_accounts, local.probe_accounts)
}

provider "aws" {
  region              = local.region
  allowed_account_ids = [local.account_id]
  default_tags {
    tags = local.tags
  }
}
