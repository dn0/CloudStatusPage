resource "aws_sqs_queue" "mon-probe" {
  name                      = "mon-probe"
  max_message_size          = 2048
  message_retention_seconds = 60

  tags = {
    cost-center = "mon-probe"
  }
}
