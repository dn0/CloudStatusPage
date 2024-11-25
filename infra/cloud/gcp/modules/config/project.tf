data "google_folder" "cloudstatus" {
  folder = "folders/123456789"
}

data "google_billing_account" "money" {
  display_name = "Money"
  open         = true
}

resource "google_project" "this" {
  project_id          = var.account_id
  name                = var.account_name
  folder_id           = data.google_folder.cloudstatus.name
  billing_account     = data.google_billing_account.money.id
  auto_create_network = false

  lifecycle {
    prevent_destroy = true
  }
}

resource "google_project_service" "this" {
  for_each = toset([
    "compute.googleapis.com",
    "logging.googleapis.com",
    "oslogin.googleapis.com",
    "pubsub.googleapis.com",
    "storage.googleapis.com",
    "storage-api.googleapis.com",
    "storage-component.googleapis.com",
  ])

  disable_on_destroy         = true
  disable_dependent_services = true
  project                    = google_project.this.id
  service                    = each.key
}
