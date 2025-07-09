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

  rds = {
    apply_immediately = {
      production = false
      testing    = true
    }[var.deployment]

    db_name    = "pg_${random_string.database_suffix.result}"
    identifier = "pg-${random_string.database_suffix.result}"

    multi_az = {
      production = false
      testing    = false
    }[var.deployment]

    monitoring_interval          = 60
    performance_insights_enabled = true

    # https://docs.aws.amazon.com/AmazonRDS/latest/PostgreSQLReleaseNotes/postgresql-versions.html.
    engine         = "postgres"
    engine_version = "16.8"

    instance_class = "db.m8g.large"

    storage_encrypted = true
    storage_type      = "gp3"

    # Storage amounts in GB for gp3 storage type.
    # This value impacts the minimal IOPS and throughput.
    # https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/CHAP_Storage.html#gp3-storage.
    allocated_storage = {
      production = 20
      testing    = 20
    }[var.deployment]

    # Config for storage autoscaling.
    max_allocated_storage = {
      production = 200
      testing    = 200
    }[var.deployment]

    # iops = {
    #   production = 3000
    #   testing    = 3000
    # }[var.deployment]

    # storage_throughput = {
    #   production = 125
    #   testing    = 125
    # }[var.deployment]

    username                    = "abodemine_admin"
    manage_master_user_password = true

    backup_retention_period = 7
    backup_window           = "03:00-04:00"
    maintenance_window      = "Sun:04:00-Sun:05:00"
    publicly_accessible     = false

    skip_final_snapshot = {
      production = false
      testing    = true
    }[var.deployment]
  }

  tags = var.tags
  vpc = {
    database_subnet_group_name = local.vpc.database_subnet_group_name
    id                         = local.vpc.vpc_id
  }
}
