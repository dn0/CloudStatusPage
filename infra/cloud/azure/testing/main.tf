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
  account    = "cloudstatus-probe-t"
  account_id = "123456"

  tags = {
    terraform   = "true"
    environment = "testing"
  }

  regions = [
    "westeurope",
    "polandcentral",
  ]
}

provider "azurerm" {
  subscription_id = local.account_id
  features {}
}

module "budget" {
  source     = "../modules/budget"
  account    = local.account
  account_id = local.account_id
}

module "config" {
  source      = "../modules/config"
  environment = local.tags.environment
  regions     = local.regions
  tags        = local.tags
}
