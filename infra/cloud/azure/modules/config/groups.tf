resource "azurerm_resource_group" "mon-agent-root" {
  name     = "mon-agent-${local.env}"
  location = "westeurope"
  tags     = var.tags
}

resource "azurerm_resource_group" "mon-agent" {
  for_each = toset(var.regions)

  name       = "mon-agent-${local.env}-${each.key}"
  location   = each.key
  managed_by = azurerm_resource_group.mon-agent-root.id
  tags = merge(var.tags, {
    cost-center = "mon-agent"
  })
}

resource "azurerm_resource_group" "mon-probe" {
  for_each = toset(var.regions)

  name       = "mon-probe-${local.env}-${each.key}"
  location   = each.key
  managed_by = azurerm_resource_group.mon-agent-root.id
  tags = merge(var.tags, {
    cost-center = "mon-probe"
  })
}
