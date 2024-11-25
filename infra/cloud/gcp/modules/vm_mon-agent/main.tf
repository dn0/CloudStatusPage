locals {
  env = substr(var.environment, 0, 1)
  labels = merge({
    cost-center = "mon-agent"
    app_name    = "mon-agent"
    app_version = var.app_version
  }, var.labels)
  startup_script_env = {
    GAR_LOCATION = var.gar_location
    ENV = merge({
      LOG_LEVEL  = "info"
      LOG_FORMAT = "text"
    }, var.env)
  }
}

data "google_compute_image" "rocky" {
  family      = "rocky-linux-9-optimized-gcp"
  project     = "rocky-linux-cloud"
  most_recent = true
}

resource "google_compute_region_instance_template" "this" {
  name         = "mon-agent-${local.env}-${var.region}-r${var.revision}"
  region       = var.region
  machine_type = var.gce_instance_type
  tags         = ["mon-agent"]
  labels       = var.labels
  metadata = {
    enable-oslogin            = "false"
    enable-oslogin-2fa        = "false"
    google-logging-enabled    = "false" # TODO: pricing?
    google-monitoring-enabled = "false" # TODO: pricing?
    startup-script            = templatefile("${path.module}/startup_script.sh", local.startup_script_env)
  }

  scheduling {
    provisioning_model          = var.gce_spot ? "SPOT" : "STANDARD"
    preemptible                 = var.gce_spot ? true : false
    automatic_restart           = var.gce_spot ? false : true
    on_host_maintenance         = "TERMINATE"
    instance_termination_action = var.gce_spot ? "STOP" : null
  }

  service_account {
    email  = var.service_account
    scopes = ["cloud-platform"]
  }

  disk {
    boot         = true
    auto_delete  = true
    device_name  = "mon-agent"
    source_image = data.google_compute_image.rocky.self_link
    disk_type    = "pd-balanced"
  }

  network_interface {
    access_config {
      network_tier = "STANDARD" # or "PREMIUM"
    }

    nic_type    = "GVNIC"
    stack_type  = "IPV4_ONLY"
    subnetwork  = var.vpc_subnet
    queue_count = 0
  }

  lifecycle {
    create_before_destroy = false
  }
}

resource "google_compute_region_health_check" "this" {
  name                = "mon-agent-${local.env}-${var.region}"
  region              = var.region
  check_interval_sec  = 3
  timeout_sec         = 1
  healthy_threshold   = 1
  unhealthy_threshold = 2

  http_health_check {
    request_path = "/healthz"
    port         = "8000"
  }
}

resource "google_compute_region_instance_group_manager" "this" {
  name               = "mon-agent-${local.env}-${var.region}-r${var.revision}"
  region             = var.region
  base_instance_name = "mon-agent"
  target_size        = 1

  version {
    name              = "mon-agent-${local.env}-r${var.revision}"
    instance_template = google_compute_region_instance_template.this.self_link
  }

  all_instances_config {
    labels = local.labels
    metadata = {
      app_version = var.app_version
    }
  }

  named_port {
    name = "http"
    port = 8000
  }

  auto_healing_policies {
    health_check      = google_compute_region_health_check.this.id
    initial_delay_sec = 180
  }

  update_policy {
    type                           = "PROACTIVE"
    instance_redistribution_type   = "NONE"
    minimal_action                 = "REPLACE"
    most_disruptive_allowed_action = "REPLACE"
    max_surge_fixed                = 5
    max_unavailable_fixed          = 5
    replacement_method             = "SUBSTITUTE"
  }

  instance_lifecycle_policy {
    force_update_on_repair    = "YES"
    default_action_on_failure = "REPAIR"
  }

  lifecycle {
    create_before_destroy = true
  }

  provisioner "local-exec" {
    when    = create
    command = "/bin/sleep 5"
  }
}
