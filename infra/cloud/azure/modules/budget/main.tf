terraform {
  required_providers {
    azurerm = {}
  }
}

resource "azurerm_cost_anomaly_alert" "this" {
  name            = "tf-cost-anomaly-alert"
  display_name    = "Daily cost anomaly by resource group"
  subscription_id = "/subscriptions/${var.account_id}"
  email_subject   = "Azure cost anomaly in ${var.account}"
  email_addresses = [var.notification_email]
}

resource "azurerm_consumption_budget_subscription" "this" {
  name            = "subscription"
  subscription_id = "/subscriptions/${var.account_id}"

  amount     = var.limit_amount
  time_grain = "Monthly"

  time_period {
    start_date = "2024-08-01T00:00:00Z"
    end_date   = "2034-08-01T00:00:00Z"
  }

  notification {
    enabled        = true
    threshold      = var.threshold
    threshold_type = "Forecasted"
    operator       = "GreaterThanOrEqualTo"
    contact_emails = [var.notification_email]
  }
}
