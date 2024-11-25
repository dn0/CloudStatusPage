variable "account" {
  description = "Azure subscription name"
  type        = string
}

variable "account_id" {
  description = "Azure subscription ID"
  type        = string
}

variable "limit_amount" {
  description = "Monthly spend limit in USD"
  type        = number
  default     = 10
}

variable "threshold" {
  description = "Threshold - percentage of the limit_amount - when the notification should be sent"
  type        = number
  default     = 90
}

variable "notification_email" {
  description = "Where to send budget alerts"
  default     = "daniel@truestatus.cloud"
}
