terraform {
  required_version = "~> 1.8.0"

  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "6.3.0"
    }
  }
}

locals {
  account_id   = "cloudstatus-probe-t"
  account_name = "cloudstatus-probe-t"

  labels = {
    terraform   = "true"
    environment = "testing"
  }

  regions = [
    "europe-west3",
  ]
}

provider "google" {
  project                         = local.account_id
  default_labels                  = local.labels
  add_terraform_attribution_label = false
}

module "config" {
  source       = "../modules/config"
  account_id   = local.account_id
  account_name = local.account_name
  regions      = local.regions
}
