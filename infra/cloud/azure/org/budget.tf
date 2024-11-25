module "budget" {
  source = "../modules/budget"

  account            = local.account
  account_id         = local.account_id
  notification_email = local.notification_email
}

resource "azurerm_consumption_budget_management_group" "root" {
  name                = "Global"
  management_group_id = data.azurerm_management_group.root.id

  amount     = 40
  time_grain = "Monthly"

  time_period {
    start_date = "2024-06-01T00:00:00Z"
    end_date   = "2034-06-01T00:00:00Z"
  }

  notification {
    enabled        = true
    threshold      = 50.0
    threshold_type = "Actual"
    operator       = "GreaterThanOrEqualTo"
    contact_emails = [local.notification_email]
  }

  notification {
    enabled        = true
    threshold      = 80.0
    threshold_type = "Forecasted"
    operator       = "GreaterThanOrEqualTo"
    contact_emails = [local.notification_email]
  }
}
