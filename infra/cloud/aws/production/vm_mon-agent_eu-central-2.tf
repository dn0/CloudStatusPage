provider "aws" {
  alias               = "eu-central-2"
  region              = "eu-central-2"
  allowed_account_ids = [local.account_id]
  default_tags { tags = local.tags }
}

module "config_eu-central-2" {
  source    = "../modules/config_region"
  providers = { aws = aws.eu-central-2 }
  account   = local.account
}

module "vm_mon-agent_eu-central-2" {
  source    = "../modules/vm_mon-agent"
  providers = { aws = aws.eu-central-2 }

  vpc_subnet_ids      = module.config_eu-central-2.vpc_subnet_ids
  ec2_role_name       = module.config_global.ec2_role_name
  s3_artifacts_bucket = module.config_global.s3_artifacts_bucket
  secret_env_arn      = module.config_global.secret_env_arn
  tags                = local.tags
  environment         = local.tags.environment
  app_version         = var.app_version
  revision            = var.app_revision
  env = merge(local.env, {
    AWS_EC2_VM_VPC_ID = module.config_eu-central-2.vpc_id,
  })

  depends_on = [module.config_eu-central-2]
}