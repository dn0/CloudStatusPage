resource "azurerm_storage_account" "mon-probe" {
  for_each = toset(var.regions)

  name                            = "probe${local.env}${each.key}"
  resource_group_name             = azurerm_resource_group.mon-probe[each.key].name
  location                        = each.key
  account_kind                    = "StorageV2"
  account_tier                    = "Standard"
  account_replication_type        = "LRS"
  access_tier                     = "Hot"
  public_network_access_enabled   = true
  allow_nested_items_to_be_public = false
  tags = merge(var.tags, {
    cost-center = "mon-probe"
  })
}

resource "azurerm_storage_container" "mon-probe" {
  for_each = toset(var.regions)

  name                  = "objects"
  storage_account_name  = azurerm_storage_account.mon-probe[each.key].name
  container_access_type = "private"
}

resource "azurerm_storage_management_policy" "mon-probe" {
  for_each = toset(var.regions)

  storage_account_id = azurerm_storage_account.mon-probe[each.key].id

  rule {
    name    = "test-cleanup"
    enabled = true
    filters {
      prefix_match = ["${azurerm_storage_container.mon-probe[each.key].name}/"]
      blob_types   = ["blockBlob"]
    }
    actions {
      base_blob {
        delete_after_days_since_modification_greater_than = 1
      }
    }
  }
}
