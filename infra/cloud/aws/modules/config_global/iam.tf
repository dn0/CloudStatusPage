resource "aws_iam_role" "ec2_mon-agent" {
  name = "ec2_mon-agent"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = "sts:AssumeRole"
        Principal = {
          Service = "ec2.amazonaws.com"
        }
      }
    ]
  })

  managed_policy_arns = [
    "arn:aws:iam::aws:policy/AmazonSSMManagedEC2InstanceDefaultPolicy",
    "arn:aws:iam::aws:policy/CloudWatchAgentServerPolicy",
  ]

  inline_policy {
    name = "inline"
    policy = jsonencode({
      Version = "2012-10-17"
      Statement = [
        {
          Effect = "Allow"
          Action = [
            "ec2:Describe*",
            "ec2:CreateTags",
          ]
          Resource = "*"
        },
        {
          Effect = "Allow"
          Action = [
            "ec2:CreateLaunchTemplateVersion",
          ]
          Resource = "*"
          Condition = {
            StringEquals = {
              "aws:ResourceTag/service" = "mon-agent",
            }
          }
        },
        {
          Sid    = "DownloadMonAgent"
          Effect = "Allow"
          Action = [
            "s3:ListMultipartUploadParts",
            "s3:ListBucketMultipartUploads",
            "s3:ListBucket",
            "s3:GetObject",
            "s3:GetBucketLocation",
          ]
          Resource = [
            data.aws_s3_bucket.artifacts.arn,
            "${data.aws_s3_bucket.artifacts.arn}/mon-agent/*",
          ]
        },
        {
          Sid    = "SecretsToEnv"
          Effect = "Allow"
          Action = [
            "secretsmanager:GetSecretValue",
          ]
          Resource = [
            "arn:aws:secretsmanager:${data.aws_region.current.name}:${var.account_id}:secret:mon-agent/*",
          ]
        },
        {
          Sid    = "BasicIAM"
          Effect = "Allow"
          Action = [
            "iam:CreateServiceLinkedRole",
            "iam:ListRoles",
            "iam:ListInstanceProfiles",
          ]
          Resource = "*"
        },
        {
          Sid    = "EC2MonitoringProbesRunInstances"
          Effect = "Allow"
          Action = [
            "ec2:RunInstances",
          ]
          Resource = [
            "arn:aws:ec2:*::image/*",
            "arn:aws:ec2:*::snapshot/*",
            "arn:aws:ec2:*:*:volume/*",
            "arn:aws:ec2:*:*:subnet/*",
            "arn:aws:ec2:*:*:network-interface/*",
            "arn:aws:ec2:*:*:security-group/*",
            "arn:aws:ec2:*:*:key-pair/*"
          ]
        },
        {
          Sid    = "EC2MonitoringProbesCreateSnapshot"
          Effect = "Allow"
          Action = [
            "ec2:CreateSnapshot",
          ]
          Resource = [
            "arn:aws:ec2:*:*:volume/*",
          ]
        },
        {
          Sid    = "EC2MonitoringProbesCreate"
          Effect = "Allow"
          Action = [
            "ec2:*",
          ]
          Resource = "*"
          Condition = {
            StringEquals = {
              "aws:RequestTag/owner" = "mon-agent",
            }
          }
        },
        {
          Sid    = "EC2MonitoringProbesModify"
          Effect = "Allow"
          Action = [
            "ec2:*",
          ]
          Resource = "*"
          Condition = {
            StringEquals = {
              "aws:ResourceTag/owner" = "mon-agent",
            }
          }
        },
        {
          Sid    = "S3MonitoringProbes"
          Effect = "Allow"
          Action = [
            "s3:*",
          ]
          Resource = [
            "arn:aws:s3:::${var.account}/",
            "arn:aws:s3:::${var.account}-*",
          ]
        },
        {
          Sid    = "SQSMonitoringProbes"
          Effect = "Allow"
          Action = [
            "sqs:*",
          ]
          Resource = [
            "arn:aws:sqs:*:*:mon-probe",
            "arn:aws:sqs:*:*:test-*",
          ]
        },
      ]
    })
  }
}
