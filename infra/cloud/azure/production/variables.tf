locals {
  env = {
    AZURE_SUBSCRIPTION_ID          = local.account_id
    AZURE_STORAGE_CONTAINER_PREFIX = "mon-probe-p-test"
    PUBSUB_PROJECT_ID              = "cloudstatus-p"

    COMMON_VPC_INTRA_PING_PONG_TIMEOUT = "40s"
    AZURE_VPC_INTER_PING_INTERVAL      = "0"
  }
}

variable "app_version" {
  type = string
}

variable "app_revision" {
  type    = number
  default = 0
}
