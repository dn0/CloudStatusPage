resource "azurerm_storage_account" "artifacts" {
  name                            = "cloudstatusartifacts"
  resource_group_name             = azurerm_resource_group.cloudstatus.name
  location                        = azurerm_resource_group.cloudstatus.location
  account_kind                    = "StorageV2"
  account_tier                    = "Standard"
  account_replication_type        = "ZRS"
  access_tier                     = "Hot"
  public_network_access_enabled   = true # NOTE: VMs will download blobs using public endpoints
  allow_nested_items_to_be_public = false
  tags                            = local.tags
}

resource "azurerm_storage_container" "mon-agent" {
  name                  = "mon-agent"
  storage_account_name  = azurerm_storage_account.artifacts.name
  container_access_type = "private"
}

resource "azurerm_role_assignment" "artifacts-admin" {
  role_definition_name = "Storage Blob Data Contributor"
  principal_id         = azuread_group.admin.id
  scope                = azurerm_storage_account.artifacts.id
}
