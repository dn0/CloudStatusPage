variable "environment" {
  description = "Deployment env"
  type        = string
}

variable "regions" {
  description = "List of enabled regions"
  type        = list(string)
}

variable "tags" {
  type    = map(string)
  default = {}
}
