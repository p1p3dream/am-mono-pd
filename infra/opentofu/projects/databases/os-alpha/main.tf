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

  os_domain = {
    domain_name = "os-${random_string.database_suffix.result}"

    # https://docs.aws.amazon.com/opensearch-service/latest/developerguide/release-notes.html.
    engine_version = {
      production = "OpenSearch_2.17"
      testing    = "OpenSearch_2.17"
    }[var.deployment]

    cluster_config = {
      instance_type = {
        production = "m7g.large.search"
        testing    = "m7g.large.search"
      }[var.deployment]

      instance_count = {
        production = 3
        testing    = 1
      }[var.deployment]

      dedicated_master_enabled = true

      dedicated_master_type = {
        production = "m7g.medium.search"
        testing    = "m7g.medium.search"
      }[var.deployment]

      dedicated_master_count = {
        production = 3
        testing    = 3
      }[var.deployment]

      warm_enabled = false

      // Set to true to multi-az.
      zone_awareness_enabled = {
        production = true
        testing    = false
      }[var.deployment]

      zone_awareness_config = {
        // Set to 3 for multi-az.
        availability_zone_count = {
          production = 3
          testing    = 0
        }[var.deployment]
      }
    }

    ebs_options = {
      ebs_enabled = true
      volume_type = "gp3"

      volume_size = {
        production = 200
        testing    = 200
      }[var.deployment]

      iops = {
        production = 6000
        testing    = 6000
      }[var.deployment]

      throughput = {
        production = 500
        testing    = 500
      }[var.deployment]
    }

    encrypt_at_rest = {
      enabled = true
    }

    node_to_node_encryption = {
      enabled = true
    }

    domain_endpoint_options = {
      enforce_https       = true
      tls_security_policy = "Policy-Min-TLS-1-2-2019-07"
    }

    advanced_security_options = {
      enabled                        = true
      internal_user_database_enabled = true
      master_user_options = {
        master_user_name     = var.os_domain.master_user_name
        master_user_password = var.os_domain.master_user_password
      }
    }
  }

  tags = var.tags
  vpc = {
    id = local.vpc.vpc_id

    database_subnets = {
      // Multi-az.
      production = local.vpc.database_subnets

      // Use subnet index 1 since we're not multi-az.
      testing = [local.vpc.database_subnets[1]]
    }[var.deployment]
  }
}
