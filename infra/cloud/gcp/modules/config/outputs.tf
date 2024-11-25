output "vpc_subnet_name" {
  value = { for i in google_compute_subnetwork.vpc_subnet : i.region => i.name }
}

output "vpc_subnet_id" {
  value = { for i in google_compute_subnetwork.vpc_subnet : i.region => i.id }
}

output "locations" {
  value = { for region, cfg in local.regions : region => cfg.multi_region }
}

output "probe_bucket" {
  value = { for region, bucket in google_storage_bucket.mon-probe : region => bucket.name }
}

output "service_account" {
  value = google_service_account.mon-agent
}
