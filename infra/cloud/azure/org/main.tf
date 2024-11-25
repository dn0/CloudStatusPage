terraform {
  required_version = "~> 1.8.0"

  required_providers {
    azurerm = {
      source  = "hashicorp/azurerm"
      version = "4.2.0"
    }
    azuread = {
      source  = "hashicorp/azuread"
      version = "2.53.1"
    }
  }
}

locals {
  account    = "cloudstatus-org"
  account_id = "345678"
  tenant_id  = "456789"

  tags = {
    terraform   = "true"
    environment = "org"
  }

  org_accounts = {
    mon-org = local.account_id
  }

  probe_accounts = {
    mon-probe-t = "123456"
    mon-probe-p = "234567"
  }

  all_accounts = merge(local.org_accounts, local.probe_accounts)

  notification_email = "daniel@truestatus.cloud"
}

provider "azurerm" {
  tenant_id       = local.tenant_id
  subscription_id = local.account_id
  features {}
}

provider "azuread" {
  tenant_id = local.tenant_id
}
