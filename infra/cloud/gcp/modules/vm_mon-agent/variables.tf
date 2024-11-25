variable "gce_instance_type" {
  type    = string
  default = "e2-small"
}

variable "gce_spot" {
  type    = bool
  default = false
}

variable "environment" {
  type = string
}

variable "region" {
  type = string
}

variable "service_account" {
  type = string
}

variable "vpc_subnet" {
  type = string
}

variable "gar_location" {
  type = string
}

variable "app_version" {
  type = string
}

variable "env" {
  type = map(string)
}

variable "revision" {
  type    = number
  default = 0
}

variable "labels" {
  type    = map(string)
  default = {}
}
