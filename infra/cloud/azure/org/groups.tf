resource "azurerm_resource_group" "cloudstatus" {
  name     = "cloudstatus-central"
  location = "westeurope"
  tags     = local.tags
}
