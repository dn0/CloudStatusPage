module "vm_mon-agent" {
  source   = "../modules/vm_mon-agent"
  for_each = toset(local.regions)

  region                  = each.key
  user_identity_id        = module.config.user_identity_id
  artifacts_container_url = module.config.artifacts_container_url
  vault_url               = module.config.vault_url
  rg_name                 = module.config.rg_agent[each.key]
  subnet_id               = module.config.subnet_id[each.key]
  tags                    = local.tags
  environment             = local.tags.environment
  app_version             = var.app_version
  revision                = var.app_revision
  env = merge(local.env, {
    CLOUD_REGION                 = each.key
    AZURE_RESOURCE_GROUP         = module.config.rg_probe[each.key]
    AZURE_COMPUTE_VM_NIC_ID      = module.config.probe_vm_nic_id[each.key]
    AZURE_COMPUTE_VM_SPOT_NIC_ID = module.config.probe_vm_spot_nic_id[each.key]
    AZURE_STORAGE_ACCOUNT_NAME   = module.config.probe_storage_name[each.key]
    AZURE_SERVICEBUS_NAMESPACE   = module.config.probe_servicebus_namespace[each.key]
  })

  depends_on = [module.config]
}
