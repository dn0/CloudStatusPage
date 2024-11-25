terraform {
  required_version = "~> 1.8.0"

  required_providers {
    azurerm = {
      source  = "hashicorp/azurerm"
      version = "4.2.0"
    }
  }
}

locals {
  account    = "cloudstatus-probe-p"
  account_id = "234567"

  tags = {
    terraform   = "true"
    environment = "production"
  }

  regions = [
    "eastus",
    "eastus2",
    "westus",
    "westus2",
    "westus3",
    "centralus",
    "westcentralus",
    "northcentralus",
  ]
}

provider "azurerm" {
  subscription_id = local.account_id
  features {}
}

module "budget" {
  source       = "../modules/budget"
  account      = local.account
  account_id   = local.account_id
  limit_amount = 200
}

module "config" {
  source      = "../modules/config"
  environment = local.tags.environment
  regions     = local.regions
  tags        = local.tags
}
