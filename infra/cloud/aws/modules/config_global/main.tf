terraform {
  required_providers {
    aws = {}
  }
}

locals {
  org_account_id = "123456"
}

data "aws_region" "current" {}

data "aws_s3_bucket" "artifacts" {
  bucket = "cloudstatus-artifacts-${data.aws_region.current.name}"
}
