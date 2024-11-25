variable "limit_amount" {
  description = "Monthly spend limit in USD"
  type        = number
  default     = 20
}

variable "threshold" {
  description = "Threshold - percentage of the limit_amount - when the notification should be sent"
  type        = number
  default     = 90
}

variable "anomaly_threshold" {
  description = "Send alert if the total impact percentage is higher that this number"
  type        = number
  default     = 50
}

variable "anomaly_min_value" {
  description = "Don't send alert if the total impact amount is lower that this amount (USD)"
  type        = number
  default     = 5
}

variable "notification_email" {
  description = "Where to send budget and anomaly alerts"
  default     = "daniel@truestatus.cloud"
}
