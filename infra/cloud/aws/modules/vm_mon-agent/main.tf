terraform {
  required_providers {
    aws = {}
  }
}

data "aws_region" "current" {}

locals {
  env    = substr(var.environment, 0, 1)
  region = data.aws_region.current.name
  tags = merge({
    cost-center = "mon-agent"
    app_name    = "mon-agent"
    app_version = var.app_version
    Name        = "mon-agent"
  }, var.tags)
  goarch = {
    x86_64 = "amd64"
    arm64  = "arm64"
  }
  app_env = merge({
    CLOUD_REGION = local.region
    LOG_LEVEL    = "info"
    LOG_FORMAT   = "text"
  }, { for k, v in var.env : k => replace(v, "#region#", local.region) })
  startup_script_env = {
    SECRET_ENV_ARN    = var.secret_env_arn
    SECRET_ENV_REGION = element(split(":", var.secret_env_arn), 3) # There is only one secret in eu-central-1
    GOARCH            = local.goarch[var.ec2_instance_arch]
    ARTIFACTS_BUCKET  = var.s3_artifacts_bucket
    ENV               = local.app_env
  }
}

data "aws_ami" "amazonlinux" {
  most_recent = true

  filter {
    name   = "owner-alias"
    values = ["amazon"]
  }

  filter {
    name   = "name"
    values = ["al2023-ami-minimal-*-${var.ec2_instance_arch}"]
  }
}

resource "aws_iam_instance_profile" "this" {
  name = "mon-agent-${local.region}"
  role = var.ec2_role_name
}

resource "aws_launch_template" "this" {
  name          = "mon-agent-${local.env}-${local.region}-r${var.revision}"
  image_id      = data.aws_ami.amazonlinux.id
  instance_type = var.ec2_instance_type
  tags          = local.tags
  network_interfaces {
    associate_public_ip_address = true
    security_groups             = var.security_group_ids
  }
  private_dns_name_options {
    enable_resource_name_dns_a_record = true
    hostname_type                     = "resource-name"
  }
  iam_instance_profile {
    name = aws_iam_instance_profile.this.name
  }
  metadata_options {
    http_endpoint          = "enabled"
    http_tokens            = "optional"
    instance_metadata_tags = "enabled"
  }
  tag_specifications {
    resource_type = "instance"
    tags          = local.tags
  }
  tag_specifications {
    resource_type = "volume"
    tags          = local.tags
  }
  user_data = base64encode(templatefile("${path.module}/startup_script.sh", local.startup_script_env))

  lifecycle {
    create_before_destroy = true
  }
}

resource "aws_autoscaling_group" "this" {
  name                = "mon-agent-${local.env}-${local.region}-r${var.revision}"
  vpc_zone_identifier = var.vpc_subnet_ids
  desired_capacity    = 1
  min_size            = 1
  max_size            = 1
  launch_template {
    name    = aws_launch_template.this.name
    version = aws_launch_template.this.latest_version
  }
  instance_refresh {
    strategy = "Rolling"
    preferences {
      instance_warmup        = 180
      checkpoint_delay       = 600
      min_healthy_percentage = 0   # 100 for less downtime
      max_healthy_percentage = 100 # 200 for less downtime
    }
    triggers = ["tag"]
  }
  dynamic "tag" {
    for_each = local.tags
    content {
      key                 = tag.key
      value               = tag.value
      propagate_at_launch = true
    }
  }

  lifecycle {
    create_before_destroy = true
  }
}
