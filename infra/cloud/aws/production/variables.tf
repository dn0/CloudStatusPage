locals {
  env = {
    AWS_S3_OBJECT_BUCKET_NAME = "${local.account}-#region#"
    AWS_S3_BUCKET_PREFIX      = "${local.account}-#region#-test"
    PUBSUB_PROJECT_ID         = "cloudstatus-p"

    COMMON_VPC_INTRA_PING_PONG_TIMEOUT = "40s"
    AWS_VPC_INTER_PING_INTERVAL        = "0"
  }
}

variable "app_version" {
  type = string
}

variable "app_revision" {
  type    = number
  default = 1
}
