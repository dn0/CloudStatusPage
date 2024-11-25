resource "google_compute_network" "vpc" {
  name                    = "vpc"
  auto_create_subnetworks = false
  routing_mode            = "REGIONAL"

  depends_on = [
    google_project_service.this
  ]
}

resource "google_compute_subnetwork" "vpc_subnet" {
  for_each = local.subnets

  name                     = each.key
  network                  = google_compute_network.vpc.name
  region                   = each.key
  ip_cidr_range            = each.value.ip_range
  private_ip_google_access = true
}

resource "google_compute_firewall" "icmp-allow-internal" {
  name          = "icmp-allow-internal"
  network       = google_compute_network.vpc.name
  priority      = 1000
  direction     = "INGRESS"
  source_ranges = ["10.0.0.0/8"]

  allow {
    protocol = "icmp"
  }
}

resource "google_compute_firewall" "http-health-check" {
  name          = "http-health-check"
  network       = google_compute_network.vpc.name
  direction     = "INGRESS"
  priority      = 1000
  target_tags   = ["mon-agent"]
  source_ranges = ["35.191.0.0/16", "130.211.0.0/22"]

  allow {
    protocol = "tcp"
    ports    = ["8000"]
  }
}

resource "google_compute_firewall" "ssh-from-iap" {
  name          = "ssh-allow"
  network       = google_compute_network.vpc.name
  direction     = "INGRESS"
  priority      = 1000
  target_tags   = ["mon-agent"]
  source_ranges = ["35.235.240.0/20"]

  allow {
    protocol = "tcp"
    ports    = ["22"]
  }
}
