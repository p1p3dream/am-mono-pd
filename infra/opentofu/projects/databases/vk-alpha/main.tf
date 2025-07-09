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
  database_id = "${var.project}-${random_string.database_suffix.result}"
  vpc         = data.terraform_remote_state.main.outputs.vpc
}

################################################################################
# MODULE: database
################################################################################

resource "random_string" "database_suffix" {
  length  = 8
  special = false
  upper   = false
}

module "database" {
  source = "./modules/database"

  aws_account_id = var.aws_account_id
  aws_region     = var.aws_region
  deployment     = var.deployment

  module_id = local.database_id

  elasticache_replication_group = {
    apply_immediately = {
      production = false
      testing    = true
    }[var.deployment]

    replication_group_id    = "vk-${random_string.database_suffix.result}"
    node_type               = "cache.m7g.large"
    num_node_groups         = 1
    replicas_per_node_group = 1
    description             = "Valkey Cluster ${random_string.database_suffix.result}"

    # https://docs.aws.amazon.com/AmazonElastiCache/latest/dg/supported-engine-versions.html.
    engine         = "valkey"
    engine_version = "8.0"
    # https://docs.aws.amazon.com/AmazonElastiCache/latest/dg/ParameterGroups.Engine.html.
    parameter_group_name = "default.valkey8"

    port = 6379

    maintenance_window = "Sun:04:00-Sun:05:00"
  }

  tags = var.tags
  vpc = {
    elasticache_subnet_group_name = local.vpc.elasticache_subnet_group_name
    id                            = local.vpc.vpc_id
  }
}
