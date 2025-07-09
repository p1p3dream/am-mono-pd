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
# MODULE: server
################################################################################

resource "random_string" "server_suffix" {
  length  = 8
  special = false
  upper   = false
}

locals {
  server_id = "${var.project.slug}-${random_string.server_suffix.result}"
}

module "server" {
  source = "./modules/server"

  aws_account_id = var.aws_account_id
  aws_region     = var.aws_region
  deployment     = var.deployment
  domains        = var.domains

  ecr_repository_names = {
    servers_go_saas = "${var.deployment}/${local.server_id}",
  }

  ecs_clusters = {
    main_fargate = {
      id = data.terraform_remote_state.servers_shared.outputs.ecs_clusters.main_fargate.id
    }
  }

  iam_roles = {
    main = {
      name = "ecs-task-main-${local.server_id}"
    }
  }

  load_balancers = {
    main = {
      listeners = {
        https = {
          arn = data.terraform_remote_state.servers_shared.outputs.load_balancers.main.listeners.https.arn
        }
      }

      security_groups = data.terraform_remote_state.servers_shared.outputs.load_balancers.main.load_balancer.security_groups
    }
  }

  module_id     = local.server_id
  module_suffix = random_string.server_suffix.result

  project = var.project
  tags    = var.tags

  vpc = {
    id              = data.terraform_remote_state.main.outputs.vpc.vpc_id
    private_subnets = data.terraform_remote_state.main.outputs.vpc.private_subnets
  }
}
