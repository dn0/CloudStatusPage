resource "google_service_account" "mon-agent" {
  account_id   = "mon-agent"
  display_name = "mon-agent-gcp"
  description  = "Service account used by mon-agent-gcp's probes"
}

resource "google_project_iam_policy" "this" {
  project = google_project.this.id
  policy_data = jsonencode({
    bindings = [
      {
        role = "roles/compute.serviceAgent"
        members = [
          "serviceAccount:service-${google_project.this.number}@compute-system.iam.gserviceaccount.com",
        ]
      },
      {
        role = "roles/editor"
        members = [
          "serviceAccount:${google_project.this.number}@cloudservices.gserviceaccount.com",
        ]
      },
      {
        role = "roles/pubsub.serviceAgent"
        members = [
          "serviceAccount:service-${google_project.this.number}@gcp-sa-pubsub.iam.gserviceaccount.com",
        ]
      },
      {
        role = "roles/compute.viewer"
        members = [
          "serviceAccount:${google_service_account.mon-agent.email}",
        ]
      },
      {
        role = "roles/compute.networkUser"
        members = [
          "serviceAccount:${google_service_account.mon-agent.email}",
        ]
      },
      {
        role = "roles/compute.admin"
        members = [
          "serviceAccount:${google_service_account.mon-agent.email}",
        ]
        condition = {
          title       = "mon-agent_probe_compute_admin"
          description = "mon-agent should be able to manage compute resources that begin with \"test-\""
          expression  = <<EOF
(resource.type == 'compute.googleapis.com/Instance' && resource.name.extract('/instances/{name}').startsWith('test-')) ||
(resource.type == 'compute.googleapis.com/Disk'     && resource.name.extract('/disks/{name}').startsWith('test-')) ||
(resource.type == 'compute.googleapis.com/Disk'     && resource.name.extract('/disks/{name}').startsWith('mon-agent-')) ||
(resource.type == 'compute.googleapis.com/Snapshot' && resource.name.extract('/snapshots/{name}').startsWith('test-'))
EOF
        }
      },
      {
        role = "roles/storage.admin"
        members = [
          "serviceAccount:${google_service_account.mon-agent.email}",
        ]
        condition = {
          title       = "mon-agent_probe_storage_admin"
          description = "mon-agent should be able to create and delete mon-probe buckets and objects"
          expression  = "resource.name.startsWith(\"projects/_/buckets/${google_project.this.project_id}\")"
        }
      },
      {
        role = "roles/pubsub.editor"
        members = [
          "serviceAccount:${google_service_account.mon-agent.email}",
        ]
      },
    ]
  })
}
