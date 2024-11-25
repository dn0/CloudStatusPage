output "s3_artifacts_bucket" {
  value = data.aws_s3_bucket.artifacts.bucket
}

output "ec2_role_name" {
  value = aws_iam_role.ec2_mon-agent.name
}

locals {
  secret_envs = {
    123456 = "arn:aws:secretsmanager:eu-central-1:123456:secret:mon-agent/env-abcEFG" # production
    234567 = "arn:aws:secretsmanager:eu-central-1:234567:secret:mon-agent/env-HIJklm" # testing
  }
}

output "secret_env_arn" {
  # value = data.aws_secretsmanager_secret.mon-agent.arn
  # All agents are reading only one secret stored in eu-central-1
  value = local.secret_envs[var.account_id]
}
