resource "azurerm_user_assigned_identity" "mon-agent" {
  name                = "mon-agent-${local.env}"
  location            = "westeurope"
  resource_group_name = azurerm_resource_group.mon-agent-root.name
}

resource "azurerm_role_assignment" "mon-agent-artifacts" {
  role_definition_name = "Storage Blob Data Reader"
  principal_id         = azurerm_user_assigned_identity.mon-agent.principal_id
  scope                = "/subscriptions/345678/resourceGroups/cloudstatus-central/providers/Microsoft.Storage/storageAccounts/cloudstatusartifacts/blobServices/default/containers/mon-agent"
}

resource "azurerm_role_assignment" "mon-agent-vm" {
  for_each = toset(var.regions)

  role_definition_name = "Virtual Machine Contributor"
  principal_id         = azurerm_user_assigned_identity.mon-agent.principal_id
  scope                = azurerm_resource_group.mon-agent[each.key].id
}

resource "azurerm_role_assignment" "mon-agent-probe-vm" {
  for_each = toset(var.regions)

  role_definition_name = "Virtual Machine Contributor"
  principal_id         = azurerm_user_assigned_identity.mon-agent.principal_id
  scope                = azurerm_resource_group.mon-probe[each.key].id
}

resource "azurerm_role_assignment" "mon-agent-net" {
  for_each = toset(var.regions)

  role_definition_name = "Network Contributor"
  principal_id         = azurerm_user_assigned_identity.mon-agent.principal_id
  scope                = azurerm_resource_group.mon-agent[each.key].id
}

resource "azurerm_role_assignment" "mon-agent-probe-net" {
  for_each = toset(var.regions)

  role_definition_name = "Network Contributor"
  principal_id         = azurerm_user_assigned_identity.mon-agent.principal_id
  scope                = azurerm_resource_group.mon-probe[each.key].id
}

resource "azurerm_role_assignment" "mon-agent-vm-disk-snapshot" {
  for_each = toset(var.regions)

  role_definition_name = "Disk Snapshot Contributor"
  principal_id         = azurerm_user_assigned_identity.mon-agent.principal_id
  scope                = azurerm_resource_group.mon-agent[each.key].id
}

resource "azurerm_role_assignment" "mon-agent-probe-vm-disk-snapshot" {
  for_each = toset(var.regions)

  role_definition_name = "Disk Snapshot Contributor"
  principal_id         = azurerm_user_assigned_identity.mon-agent.principal_id
  scope                = azurerm_resource_group.mon-probe[each.key].id
}

resource "azurerm_role_assignment" "mon-agent-probe-storage" {
  for_each = toset(var.regions)

  role_definition_name = "Storage Blob Data Contributor"
  principal_id         = azurerm_user_assigned_identity.mon-agent.principal_id
  scope                = azurerm_storage_account.mon-probe[each.key].id
}

resource "azurerm_role_assignment" "admin-probe-storage" {
  for_each = toset(var.regions)

  role_definition_name = "Storage Blob Data Contributor"
  principal_id         = "abcdef" # admin group
  scope                = azurerm_storage_account.mon-probe[each.key].id
}

resource "azurerm_role_assignment" "mon-agent-probe-servicebus" {
  for_each = toset(var.regions)

  role_definition_name = "Azure Service Bus Data Owner"
  principal_id         = azurerm_user_assigned_identity.mon-agent.principal_id
  scope                = azurerm_servicebus_namespace.mon-probe[each.key].id
}

resource "azurerm_role_assignment" "admin-probe-servicebus" {
  for_each = toset(var.regions)

  role_definition_name = "Azure Service Bus Data Owner"
  principal_id         = "abcdef" # admin group
  scope                = azurerm_servicebus_namespace.mon-probe[each.key].id
}

resource "azurerm_role_assignment" "mon-agent-key-vault" {
  role_definition_name = "Key Vault Secrets User"
  principal_id         = azurerm_user_assigned_identity.mon-agent.principal_id
  scope                = azurerm_key_vault.mon-probe.id
}

resource "azurerm_role_assignment" "admin-key-vault" {
  role_definition_name = "Key Vault Administrator"
  principal_id         = "abcdef" # admin group
  scope                = azurerm_key_vault.mon-probe.id
}
