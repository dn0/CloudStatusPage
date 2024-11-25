terraform {
  required_providers {
    aws = {}
  }
}

resource "aws_budgets_budget" "this" {
  name         = "account"
  budget_type  = "COST"
  limit_amount = var.limit_amount
  limit_unit   = "USD"
  time_unit    = "MONTHLY"

  notification {
    comparison_operator        = "GREATER_THAN"
    threshold                  = var.threshold
    threshold_type             = "PERCENTAGE"
    notification_type          = "FORECASTED"
    subscriber_email_addresses = [var.notification_email]
  }
}

resource "aws_ce_anomaly_monitor" "spend_per_service" {
  name              = "costs_per_service"
  monitor_type      = "DIMENSIONAL"
  monitor_dimension = "SERVICE"
}

resource "aws_ce_anomaly_subscription" "spend_per_service" {
  name      = "DAILYSUBSCRIPTION"
  frequency = "DAILY" # "IMMEDIATE" requires SNS

  monitor_arn_list = [
    aws_ce_anomaly_monitor.spend_per_service.arn,
  ]

  subscriber {
    type    = "EMAIL"
    address = var.notification_email
  }

  threshold_expression {
    and {
      dimension {
        key           = "ANOMALY_TOTAL_IMPACT_PERCENTAGE"
        match_options = ["GREATER_THAN_OR_EQUAL"]
        values        = [var.anomaly_threshold]
      }
    }
    and {
      dimension {
        key           = "ANOMALY_TOTAL_IMPACT_ABSOLUTE"
        match_options = ["GREATER_THAN_OR_EQUAL"]
        values        = ["${var.anomaly_min_value}"] # Ignore anomalies with amount lower than this
      }
    }
  }
}
