module "vm_mon-agent" {
  source   = "../modules/vm_mon-agent"
  for_each = toset(local.regions)

  region          = each.key
  gce_spot        = true
  service_account = module.config.service_account.email
  vpc_subnet      = module.config.vpc_subnet_name[each.key]
  gar_location    = module.config.locations[each.key]
  labels          = local.labels
  environment     = local.labels.environment
  app_version     = var.app_version
  revision        = var.app_revision
  env = merge(local.env, {
    CLOUD_REGION                     = each.key
    GCP_COMPUTE_VM_SUBNETWORK        = module.config.vpc_subnet_id[each.key]
    GCP_COMPUTE_VM_PREFIX            = "test-${each.key}"
    GCP_COMPUTE_VM_SPOT_PREFIX       = "test-spot-${each.key}"
    GCP_COMPUTE_DISK_SNAPSHOT_PREFIX = "test-${each.key}"
    GCP_STORAGE_OBJECT_BUCKET_NAME   = module.config.probe_bucket[each.key]
    GCP_STORAGE_BUCKET_PREFIX        = "${module.config.probe_bucket[each.key]}-test"
    GCP_PUBSUB_PROJECT               = local.account_id
    GCP_PUBSUB_TOPIC                 = "mon-probe-${each.key}"
    GCP_PUBSUB_SUBSCRIPTION          = "mon-probe-${each.key}"
    GCP_COMPUTE_VM_ZONE_SKIP         = "europe-west9-a"
  })

  depends_on = [module.config]
}
