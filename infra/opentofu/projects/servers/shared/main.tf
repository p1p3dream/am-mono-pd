terraform {
  backend "s3" {
    region         = var.aws_region
    bucket         = var.s3_backend_bucket
    key            = var.s3_backend_key
    dynamodb_table = var.s3_backend_table
  }

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.91.0"
    }

    random = {
      source  = "hashicorp/random"
      version = "~> 3.6.3"
    }
  }

  required_version = "~> 1.7"
}

locals {
  server_id = "${var.project}-${random_string.server_suffix.result}"
  vpc       = data.terraform_remote_state.main.outputs.vpc
}

################################################################################
# MODULE: server
################################################################################

resource "random_string" "server_suffix" {
  length  = 8
  special = false
  upper   = false
}

module "server" {
  source = "./modules/server"

  aws_account_id = var.aws_account_id
  aws_region     = var.aws_region
  deployment     = var.deployment
  domains        = var.domains

  ecs_clusters = {
    main_fargate = {
      name = "main-fargate-${local.server_id}"
      setting = {
        name  = "containerInsights"
        value = "enhanced"
      }
    }
  }

  management_aws_profile = var.management_aws_profile
  module_id              = local.server_id

  tags = var.tags
  vpc = {
    id             = local.vpc.vpc_id
    public_subnets = local.vpc.public_subnets
  }
}
