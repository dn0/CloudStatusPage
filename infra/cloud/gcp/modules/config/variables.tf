variable "account_id" {
  type = string
}

variable "account_name" {
  type = string
}

variable "regions" {
  description = "List of enabled regions"
  type        = list(string)
}
