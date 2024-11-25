output "user_identity_id" {
  value = azurerm_user_assigned_identity.mon-agent.id
}

output "artifacts_container_url" {
  value = "https://cloudstatusartifacts.blob.core.windows.net/mon-agent"
}

output "vault_url" {
  value = azurerm_key_vault.mon-probe.vault_uri
}

output "rg_agent" {
  value = { for i in azurerm_resource_group.mon-agent : i.location => i.name }
}

output "rg_probe" {
  value = { for i in azurerm_resource_group.mon-probe : i.location => i.name }
}

output "probe_storage_name" {
  value = { for k, v in azurerm_storage_account.mon-probe : k => v.name }
}

output "probe_container_name" {
  value = { for k, v in azurerm_storage_container.mon-probe : k => v.name }
}

output "subnet_id" {
  value = { for k, v in azurerm_subnet.mon-agent : k => v.id }
}

output "probe_vm_nic_id" {
  value = { for k, v in azurerm_network_interface.mon-probe-vm : k => v.id }
}

output "probe_vm_spot_nic_id" {
  value = { for k, v in azurerm_network_interface.mon-probe-vm-spot : k => v.id }
}

output "probe_servicebus_namespace" {
  value = { for k, v in azurerm_servicebus_namespace.mon-probe : k => v.name }
}
