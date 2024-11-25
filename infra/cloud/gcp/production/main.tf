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
  account_id   = "cloudstatus-probe-p"
  account_name = "cloudstatus-probe-p"

  labels = {
    terraform   = "true"
    environment = "production"
  }

  regions = [
    "us-west1",
    "us-east1",
    "us-central1",
    "us-east4",
    "us-east5",
    "us-south1",
    "us-west2",
    "us-west3",
    "us-west4",
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
