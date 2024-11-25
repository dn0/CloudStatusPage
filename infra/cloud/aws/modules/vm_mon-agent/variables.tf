variable "ec2_instance_arch" {
  type    = string
  default = "arm64" # arm64 or x86_64
}

variable "ec2_instance_type" {
  type    = string
  default = "t4g.nano"
}

variable "ec2_role_name" {
  type = string
}

variable "environment" {
  type = string
}

variable "security_group_ids" {
  type    = list(string)
  default = []
}

variable "vpc_subnet_ids" {
  type = list(string)
}

variable "s3_artifacts_bucket" {
  type = string
}

variable "app_version" {
  type = string
}

variable "env" {
  type = map(string)
}

variable "secret_env_arn" {
  type = string
}

variable "revision" {
  type    = number
  default = 0
}

variable "tags" {
  type    = map(string)
  default = {}
}
