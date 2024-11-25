terraform {
  required_providers {
    aws = {}
  }
}

data "aws_region" "current" {}

locals {
  region = data.aws_region.current.name

  regions = {
    af-south-1 = {
      cidr           = "10.11.100.0/23"
      azs            = ["af-south-1a", "af-south-1b", "af-south-1c"]
      public_subnets = ["10.11.100.0/27", "10.11.100.32/27", "10.11.100.64/27"]
    }

    ap-east-1 = {
      cidr           = "10.12.100.0/23"
      azs            = ["ap-east-1a", "ap-east-1b", "ap-east-1c"]
      public_subnets = ["10.12.100.0/27", "10.12.100.32/27", "10.12.100.64/27"]
    }

    ap-northeast-1 = {
      cidr           = "10.13.100.0/23"
      azs            = ["ap-northeast-1a", "ap-northeast-1c", "ap-northeast-1d"]
      public_subnets = ["10.13.100.0/27", "10.13.100.32/27", "10.13.100.64/27"]
    }

    ap-northeast-2 = {
      cidr           = "10.14.100.0/23"
      azs            = ["ap-northeast-2a", "ap-northeast-2b", "ap-northeast-2c", "ap-northeast-2d"]
      public_subnets = ["10.14.100.0/27", "10.14.100.32/27", "10.14.100.64/27", "10.14.100.96/27"]
    }

    ap-northeast-3 = {
      cidr           = "10.15.100.0/23"
      azs            = ["ap-northeast-3a", "ap-northeast-3b", "ap-northeast-3c"]
      public_subnets = ["10.15.100.0/27", "10.15.100.32/27", "10.15.100.64/27"]
    }

    ap-south-1 = {
      cidr           = "10.16.100.0/23"
      azs            = ["ap-south-1a", "ap-south-1b", "ap-south-1c"]
      public_subnets = ["10.16.100.0/27", "10.16.100.32/27", "10.16.100.64/27"]
    }

    ap-southeast-1 = {
      cidr           = "10.17.100.0/23"
      azs            = ["ap-southeast-1a", "ap-southeast-1b", "ap-southeast-1c"]
      public_subnets = ["10.17.100.0/27", "10.17.100.32/27", "10.17.100.64/27"]
    }

    ap-southeast-2 = {
      cidr           = "10.18.100.0/23"
      azs            = ["ap-southeast-2a", "ap-southeast-2b", "ap-southeast-2c"]
      public_subnets = ["10.18.100.0/27", "10.18.100.32/27", "10.18.100.64/27"]
    }

    ap-southeast-3 = {
      cidr           = "10.19.100.0/23"
      azs            = ["ap-southeast-3a", "ap-southeast-3b", "ap-southeast-3c"]
      public_subnets = ["10.19.100.0/27", "10.19.100.32/27", "10.19.100.64/27"]
    }

    ap-southeast-4 = {
      cidr           = "10.20.100.0/23"
      azs            = ["ap-southeast-4a", "ap-southeast-4b", "ap-southeast-4c"]
      public_subnets = ["10.20.100.0/27", "10.20.100.32/27", "10.20.100.64/27"]
    }

    ap-southeast-5 = {
      cidr           = "10.21.100.0/23"
      azs            = ["ap-southeast-5a", "ap-southeast-5b", "ap-southeast-5c"]
      public_subnets = ["10.21.100.0/27", "10.21.100.32/27", "10.21.100.64/27"]
    }

    ca-central-1 = {
      cidr           = "10.22.100.0/23"
      azs            = ["ca-central-1a", "ca-central-1b", "ca-central-1d"]
      public_subnets = ["10.22.100.0/27", "10.22.100.32/27", "10.22.100.64/27"]
    }

    ca-west-1 = {
      cidr           = "10.23.100.0/23"
      azs            = ["ca-west-1a", "ca-west-1b", "ca-west-1c"]
      public_subnets = ["10.23.100.0/27", "10.23.100.32/27", "10.23.100.64/27"]
    }

    eu-central-1 = {
      cidr           = "10.24.100.0/23"
      azs            = ["eu-central-1a", "eu-central-1b", "eu-central-1c"]
      public_subnets = ["10.24.100.0/27", "10.24.100.32/27", "10.24.100.64/27"]
    }

    eu-central-2 = {
      cidr           = "10.25.100.0/23"
      azs            = ["eu-central-2a", "eu-central-2b", "eu-central-2c"]
      public_subnets = ["10.25.100.0/27", "10.25.100.32/27", "10.25.100.64/27"]
    }

    eu-north-1 = {
      cidr           = "10.26.100.0/23"
      azs            = ["eu-north-1a", "eu-north-1b", "eu-north-1c"]
      public_subnets = ["10.26.100.0/27", "10.26.100.32/27", "10.26.100.64/27"]
    }

    eu-south-1 = {
      cidr           = "10.27.100.0/23"
      azs            = ["eu-south-1a", "eu-south-1b", "eu-south-1c"]
      public_subnets = ["10.27.100.0/27", "10.27.100.32/27", "10.27.100.64/27"]
    }

    eu-south-2 = {
      cidr           = "10.28.100.0/23"
      azs            = ["eu-south-2a", "eu-south-2b", "eu-south-2c"]
      public_subnets = ["10.28.100.0/27", "10.28.100.32/27", "10.28.100.64/27"]
    }

    eu-west-1 = {
      cidr           = "10.29.100.0/23"
      azs            = ["eu-west-1a", "eu-west-1b", "eu-west-1c"]
      public_subnets = ["10.29.100.0/27", "10.29.100.32/27", "10.29.100.64/27"]
    }

    eu-west-2 = {
      cidr           = "10.30.100.0/23"
      azs            = ["eu-west-2a", "eu-west-2b", "eu-west-2c"]
      public_subnets = ["10.30.100.0/27", "10.30.100.32/27", "10.30.100.64/27"]
    }

    eu-west-3 = {
      cidr           = "10.31.100.0/23"
      azs            = ["eu-west-3a", "eu-west-3b", "eu-west-3c"]
      public_subnets = ["10.31.100.0/27", "10.31.100.32/27", "10.31.100.64/27"]
    }

    il-central-1 = {
      cidr           = "10.32.100.0/23"
      azs            = ["il-central-1a", "il-central-1b", "il-central-1c"]
      public_subnets = ["10.32.100.0/27", "10.32.100.32/27", "10.32.100.64/27"]
    }

    me-central-1 = {
      cidr           = "10.33.100.0/23"
      azs            = ["me-central-1a", "me-central-1b", "me-central-1c"]
      public_subnets = ["10.33.100.0/27", "10.33.100.32/27", "10.33.100.64/27"]
    }

    me-south-1 = {
      cidr           = "10.34.100.0/23"
      azs            = ["me-south-1a", "me-south-1b", "me-south-1c"]
      public_subnets = ["10.34.100.0/27", "10.34.100.32/27", "10.34.100.64/27"]
    }

    sa-east-1 = {
      cidr           = "10.35.100.0/23"
      azs            = ["sa-east-1a", "sa-east-1b", "sa-east-1c"]
      public_subnets = ["10.35.100.0/27", "10.35.100.32/27", "10.35.100.64/27"]
    }

    us-east-1 = {
      cidr           = "10.36.100.0/23"
      azs            = ["us-east-1a", "us-east-1b", "us-east-1c", "us-east-1d", "us-east-1e", "us-east-1f"]
      public_subnets = ["10.36.100.0/27", "10.36.100.32/27", "10.36.100.64/27", "10.36.100.96/27", "10.36.100.128/27", "10.36.100.160/27"]
    }

    us-east-2 = {
      cidr           = "10.37.100.0/23"
      azs            = ["us-east-2a", "us-east-2b", "us-east-2c"]
      public_subnets = ["10.37.100.0/27", "10.37.100.32/27", "10.37.100.64/27"]
    }

    us-west-1 = {
      cidr           = "10.38.100.0/23"
      azs            = ["us-west-1a", "us-west-1b"]
      public_subnets = ["10.38.100.0/27", "10.38.100.32/27"]
    }

    us-west-2 = {
      cidr           = "10.39.100.0/23"
      azs            = ["us-west-2a", "us-west-2b", "us-west-2c", "us-west-2d"]
      public_subnets = ["10.39.100.0/27", "10.39.100.32/27", "10.39.100.64/27", "10.39.100.96/27"]
    }
  }

  region_cfg = local.regions[local.region]
}
