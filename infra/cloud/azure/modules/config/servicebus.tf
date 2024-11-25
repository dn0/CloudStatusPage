resource "azurerm_servicebus_namespace" "mon-probe" {
  for_each = toset(var.regions)

  name                          = "probe${local.env}${each.key}"
  resource_group_name           = azurerm_resource_group.mon-probe[each.key].name
  location                      = each.key
  sku                           = "Basic"
  public_network_access_enabled = true
  tags = merge(var.tags, {
    cost-center = "mon-probe"
  })
}

resource "azurerm_servicebus_queue" "mon-probe" {
  for_each = toset(var.regions)

  name                = "mon-probe"
  namespace_id        = azurerm_servicebus_namespace.mon-probe[each.key].id
  lock_duration       = "PT1M"
  default_message_ttl = "PT1M"
}
