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
# MODULE: distribution
################################################################################

resource "random_string" "distribution_suffix" {
  length  = 8
  special = false
  upper   = false
}

locals {
  distribution_id = "${var.project.slug}-${random_string.distribution_suffix.result}"
}

module "distribution" {
  source = "./modules/distribution"

  aws_account_id = var.aws_account_id
  aws_region     = var.aws_region

  certs = {
    abodemine_main = {
      arn = data.terraform_remote_state.servers_shared.outputs.certs.abodemine_main.cert.arn
    }
  }

  deployment = var.deployment
  domains    = var.domains

  load_balancers = {
    main = {
      dns_name = data.terraform_remote_state.servers_shared.outputs.load_balancers.main.load_balancer.dns_name
    }
  }

  management_aws_profile = var.management_aws_profile
  module_id              = local.distribution_id
  module_suffix          = random_string.distribution_suffix.result

  tags = var.tags
}
