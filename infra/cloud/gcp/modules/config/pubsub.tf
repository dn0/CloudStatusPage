resource "google_pubsub_topic" "mon-probe" {
  for_each = toset(var.regions)

  name = "mon-probe-${each.key}"
  labels = {
    cost-center = "mon-probe"
    publisher   = "mon-agent"
  }

  message_retention_duration = "600s"
}

resource "google_pubsub_subscription" "mon-probe" {
  for_each = toset(var.regions)

  name  = "mon-probe-${each.key}"
  topic = google_pubsub_topic.mon-probe[each.key].id
  labels = {
    cost-center = "mon-probe"
    consumer    = "mon-agent"
  }

  message_retention_duration   = "600s"
  retain_acked_messages        = false
  ack_deadline_seconds         = 10
  enable_message_ordering      = false
  enable_exactly_once_delivery = false
}
