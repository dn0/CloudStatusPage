resource "google_storage_bucket" "mon-probe" {
  for_each = toset(var.regions)

  name          = "${var.account_id}-${each.key}"
  storage_class = "STANDARD"
  location      = upper(each.key)

  uniform_bucket_level_access = true

  labels = {
    cost-center = "mon-probe"
  }

  lifecycle_rule {
    condition {
      age = 1
    }
    action {
      type = "Delete"
    }
  }
}
