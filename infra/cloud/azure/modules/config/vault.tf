resource "azurerm_key_vault" "mon-probe" {
  name                          = "mon-probe-${local.env}"
  resource_group_name           = azurerm_resource_group.mon-agent-root.name
  location                      = "westeurope"
  tenant_id                     = data.azurerm_client_config.current.tenant_id
  sku_name                      = "standard"
  soft_delete_retention_days    = 7
  purge_protection_enabled      = false
  public_network_access_enabled = true
  enable_rbac_authorization     = true
  tags = merge(var.tags, {
    cost-center = "undefined"
  })
}
