terraform {
  backend "s3" {
    region         = var.aws_region
    bucket         = var.s3_backend_bucket
    key            = var.s3_backend_keys.task
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
# MODULE: task
################################################################################

resource "random_string" "task_suffix" {
  length  = 8
  special = false
  upper   = false
}

locals {
  task_id = "${var.task.slug}-${random_string.task_suffix.result}"
}

module "task" {
  source = "./modules/task"

  aws_account_id  = var.aws_account_id
  aws_region      = var.aws_region
  deployment      = var.deployment
  dynamodb_tables = var.dynamodb_tables

  ecr_repository_names = {
    main = "${var.deployment}/${local.task_id}",
  }

  iam_roles = {
    main = {
      name = local.task_id
    }
  }

  module_id     = local.task_id
  module_suffix = random_string.task_suffix.result

  s3_buckets = {
    partner_data = [for bucket in data.terraform_remote_state.workers_go_datapipe.outputs.s3_buckets.partner_data : bucket.arn]
  }

  security_groups = {
    main = {
      name = local.task_id
    }
  }

  tags = var.tags
  task = var.task

  vpc = {
    id              = data.terraform_remote_state.main.outputs.vpc.vpc_id
    private_subnets = data.terraform_remote_state.main.outputs.vpc.private_subnets
  }
}
