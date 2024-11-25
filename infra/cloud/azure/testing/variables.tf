locals {
  env = {
    AZURE_SUBSCRIPTION_ID          = local.account_id
    AZURE_STORAGE_CONTAINER_PREFIX = "mon-probe-t-test"

    COMMON_VPC_INTRA_PING_PONG_TIMEOUT = "40s"

    PUBSUB_PROJECT_ID                         = "" # Disable pubsub messages "cloudstatus-t"
    LOG_LEVEL                                 = "debug"
    PING_INTERVAL                             = "600s"
    PROBE_INTERVAL_DEFAULT                    = "0"
    PROBE_LONG_INTERVAL_DEFAULT               = "0"
    AZURE_STORAGE_CONTAINER_PROBE_INTERVAL    = "0"
    AZURE_STORAGE_BLOB_PROBE_INTERVAL         = "0"
    AZURE_COMPUTE_VM_PROBE_INTERVAL           = "0"
    AZURE_COMPUTE_VM_SPOT_PROBE_INTERVAL      = "0"
    AZURE_COMPUTE_VM_METADATA_PROBE_INTERVAL  = "0"
    AZURE_COMPUTE_VHD_SNAPSHOT_PROBE_INTERVAL = "0"
    AZURE_SERVICEBUS_QUEUE_MESSAGE_INTERVAL   = "0"
    AZURE_VPC_INTER_PING_INTERVAL             = "0"
  }
}

variable "app_version" {
  type = string
}

variable "app_revision" {
  type    = number
  default = 0
}
