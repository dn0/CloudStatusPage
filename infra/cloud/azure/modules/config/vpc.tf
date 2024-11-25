resource "azurerm_network_security_group" "mon-agent" {
  for_each = local.subnets

  name                = "mon-agent-${local.env}"
  location            = each.key
  resource_group_name = azurerm_resource_group.mon-agent[each.key].name

  security_rule {
    name                       = "ICMP"
    priority                   = 1000
    direction                  = "Inbound"
    access                     = "Allow"
    protocol                   = "Icmp"
    destination_port_range     = "*"
    destination_address_prefix = "*"
    source_port_range          = "*"
    source_address_prefixes = [
      "10.0.0.0/8",
    ]
  }

  security_rule {
    name                       = "SSH"
    priority                   = 1001
    direction                  = "Inbound"
    access                     = "Allow"
    protocol                   = "Tcp"
    destination_port_range     = "22"
    destination_address_prefix = "*"
    source_port_range          = "*"
    source_address_prefixes = []
  }
}

resource "azurerm_virtual_network" "vnet" {
  for_each = local.subnets

  name                = "vnet-${local.env}-${each.key}"
  location            = each.key
  resource_group_name = azurerm_resource_group.mon-agent[each.key].name
  address_space       = [each.value.ip_range]
  tags = merge(var.tags, {
    cost-center = "undefined"
  })
}

resource "azurerm_subnet" "mon-agent" {
  for_each = local.subnets

  name                            = "mon-agent-${local.env}"
  resource_group_name             = azurerm_resource_group.mon-agent[each.key].name
  virtual_network_name            = azurerm_virtual_network.vnet[each.key].name
  address_prefixes                = ["${trimsuffix(each.value.ip_range, "/24")}/28"]
  default_outbound_access_enabled = true
  service_endpoints = [
    "Microsoft.KeyVault",
    "Microsoft.Storage",
  ]
}

resource "azurerm_subnet_network_security_group_association" "mon-agent" {
  for_each = local.subnets

  subnet_id                 = azurerm_subnet.mon-agent[each.key].id
  network_security_group_id = azurerm_network_security_group.mon-agent[each.key].id
}

resource "azurerm_network_interface" "mon-probe-vm" {
  for_each = local.subnets

  name                = "test-nic"
  location            = each.key
  resource_group_name = azurerm_resource_group.mon-probe[each.key].name

  ip_configuration {
    name                          = "test-ip"
    subnet_id                     = azurerm_subnet.mon-agent[each.key].id
    private_ip_address_allocation = "Dynamic"
  }
}

resource "azurerm_network_interface" "mon-probe-vm-spot" {
  for_each = local.subnets

  name                = "test-spot-nic"
  location            = each.key
  resource_group_name = azurerm_resource_group.mon-probe[each.key].name

  ip_configuration {
    name                          = "test-spot-ip"
    subnet_id                     = azurerm_subnet.mon-agent[each.key].id
    private_ip_address_allocation = "Dynamic"
  }
}

resource "azurerm_virtual_network_peering" "peering" {
  for_each = local.regions_cartesian

  name                         = "peering-${local.env}-${each.key}"
  resource_group_name          = azurerm_resource_group.mon-agent[each.value.src].name
  virtual_network_name         = azurerm_virtual_network.vnet[each.value.src].name
  remote_virtual_network_id    = azurerm_virtual_network.vnet[each.value.dst].id
  allow_virtual_network_access = true
  allow_forwarded_traffic      = true
  allow_gateway_transit        = false
}
