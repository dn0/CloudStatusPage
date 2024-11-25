locals {
  env = {
    GCP_PROJECT_ID    = local.account_id
    PUBSUB_PROJECT_ID = "cloudstatus-p"

    COMMON_VPC_INTRA_PING_PONG_TIMEOUT = "40s"
    GCP_VPC_INTER_PING_INTERVAL        = "0"
  }
}

variable "app_version" {
  type = string
}

variable "app_revision" {
  type    = number
  default = 5
}
