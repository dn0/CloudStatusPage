locals {
  env = {
    GCP_PROJECT_ID = local.account_id

    COMMON_VPC_INTRA_PING_PONG_TIMEOUT = "40s"

    PUBSUB_PROJECT_ID                        = "" # Disable pubsub messages "cloudstatus-t"
    LOG_LEVEL                                = "debug"
    PING_INTERVAL                            = "600s"
    PROBE_INTERVAL_DEFAULT                   = "0"
    PROBE_LONG_INTERVAL_DEFAULT              = "0"
    GCP_STORAGE_BUCKET_PROBE_INTERVAL        = "0"
    GCP_STORAGE_OBJECT_PROBE_INTERVAL        = "0"
    GCP_COMPUTE_VM_PROBE_INTERVAL            = "0"
    GCP_COMPUTE_VM_SPOT_PROBE_INTERVAL       = "0"
    GCP_COMPUTE_VM_METADATA_PROBE_INTERVAL   = "0"
    GCP_COMPUTE_DISK_SNAPSHOT_PROBE_INTERVAL = "0"
    GCP_PUBSUB_MESSAGE_INTERVAL              = "0"
    GCP_VPC_INTER_PING_INTERVAL              = "0"
  }
}

variable "app_version" {
  type = string
}

variable "app_revision" {
  type    = number
  default = 0
}
