locals {
  env = {
    AWS_S3_OBJECT_BUCKET_NAME = "${local.account}-#region#"
    AWS_S3_BUCKET_PREFIX      = "${local.account}-#region#-test"
    PUBSUB_PROJECT_ID         = "" # Disable pubsub messages "cloudstatus-core-t"

    COMMON_VPC_INTRA_PING_PONG_TIMEOUT = "40s"

    LOG_LEVEL                           = "debug"
    PING_INTERVAL                       = "600s"
    PROBE_INTERVAL_DEFAULT              = "0"
    PROBE_LONG_INTERVAL_DEFAULT         = "0"
    AWS_S3_BUCKET_PROBE_INTERVAL        = "0"
    AWS_S3_OBJECT_PROBE_INTERVAL        = "0"
    AWS_EC2_VM_PROBE_INTERVAL           = "0"
    AWS_EC2_VM_SPOT_PROBE_INTERVAL      = "0"
    AWS_EC2_VM_METADATA_PROBE_INTERVAL  = "0"
    AWS_EC2_EBS_SNAPSHOT_PROBE_INTERVAL = "0"
    AWS_SQS_MESSAGE_INTERVAL            = "0"
    AWS_VPC_INTER_PING_INTERVAL         = "0"
  }
}

variable "app_version" {
  type = string
}

variable "app_revision" {
  type    = number
  default = 0
}
