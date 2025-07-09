terraform {
  backend "s3" {
    region         = var.aws_region
    bucket         = var.s3_backend_bucket
    key            = var.s3_backend_keys.project
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
  worker_id = "${var.project.slug}-${random_string.worker_suffix.result}"
}

################################################################################
# MODULE: worker
################################################################################

resource "random_string" "worker_suffix" {
  length  = 8
  special = false
  upper   = false
}

module "worker" {
  source = "./modules/worker"

  aws_region = var.aws_region
  deployment = var.deployment

  ecs_clusters = {
    main_fargate = {
      name = "main-fargate-${local.worker_id}"
      setting = {
        name  = "containerInsights"
        value = "enhanced"
      }
    }
  }

  module_id = local.worker_id
  tags      = var.tags
}
