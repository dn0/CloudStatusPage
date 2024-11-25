# NOTE: the billing account is managed manually
data "azurerm_management_group" "root" {
  display_name = "Tenant Root Group"
}

resource "azurerm_role_assignment" "root_admin" {
  role_definition_name = "Owner"
  principal_id         = azuread_group.admin.id
  scope                = data.azurerm_management_group.root.id
}

# NOTE: had to use 'az account alias delete' to make this work
resource "azurerm_subscription" "this" {
  for_each = local.all_accounts

  alias             = each.key
  subscription_name = "cloudstatus-${each.key}"
  subscription_id   = each.value
  tags              = local.tags
}

resource "azurerm_management_group_subscription_association" "root" {
  for_each = azurerm_subscription.this

  management_group_id = data.azurerm_management_group.root.id
  subscription_id     = "/subscriptions/${each.value.subscription_id}"
}
