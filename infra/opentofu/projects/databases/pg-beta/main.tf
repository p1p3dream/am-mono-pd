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

  aurora_cluster = {
    apply_immediately = {
      production = false
      testing    = true
    }[var.deployment]
    backup_retention_period = 7
    cluster_identifier      = "pg-${random_string.database_suffix.result}"
    database_name           = "pg_${random_string.database_suffix.result}"
    engine                  = "aurora-postgresql"
    # https://docs.aws.amazon.com/AmazonRDS/latest/AuroraUserGuide/USER_UpgradeDBInstance.PostgreSQL.UpgradeVersion.html.
    engine_version               = "16.6"
    master_username              = "abodemine_admin"
    performance_insights_enabled = true
    preferred_backup_window      = "03:00-04:00"

    skip_final_snapshot = {
      production = false
      testing    = true
    }[var.deployment]
  }

  aurora_instance = {
    apply_immediately = {
      production = false
      testing    = true
    }[var.deployment]

    count = {
      // 1 writer, 1 reader.
      production = 2

      // 1 writer.
      testing = 1
    }[var.deployment]

    # https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/Concepts.DBInstanceClass.Support.html.
    # https://aws.amazon.com/ec2/instance-types/r8g/.
    instance_class = {
      production = "db.r8g.2xlarge"
      testing    = "db.r8g.2xlarge"
    }[var.deployment]

    publicly_accessible = {
      production = false
      testing    = true
    }[var.deployment]
  }

  aurora_param_group = {
    # aws rds describe-db-engine-versions --engine aurora-postgresql --query "DBEngineVersions[].DBParameterGroupFamily" | jq -r '.[]' | sort | uniq
    family = "aurora-postgresql16"
  }

  tags = var.tags
  vpc = {
    azs                          = local.vpc.azs
    database_subnets_cidr_blocks = local.vpc.database_subnets_cidr_blocks
    database_subnet_group_name   = local.vpc.database_subnet_group_name
    id                           = local.vpc.vpc_id
    private_subnets_cidr_blocks  = local.vpc.private_subnets_cidr_blocks
  }
}
