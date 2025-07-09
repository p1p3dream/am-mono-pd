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

################################################################################
# MODULE: worker
################################################################################

resource "random_string" "worker_suffix" {
  length  = 8
  special = false
  upper   = false
}

locals {
  worker_id = "${var.project.slug}-${random_string.worker_suffix.result}"
}

module "worker" {
  source = "./modules/worker"

  aws_account_id = var.aws_account_id
  aws_region     = var.aws_region
  deployment     = var.deployment

  dynamodb_tables = {
    secure_download = {
      name = "secure-download-${random_string.worker_suffix.result}"
    }
  }

  iam_roles = {}

  module_id     = local.worker_id
  module_suffix = random_string.worker_suffix.result
  project       = var.project

  s3_buckets = {
    secure_download = {
      name = "secure-download-${random_string.worker_suffix.result}"
    }
  }

  tags = var.tags

  vpc = {
    id              = data.terraform_remote_state.main.outputs.vpc.vpc_id
    private_subnets = data.terraform_remote_state.main.outputs.vpc.private_subnets
  }
}
