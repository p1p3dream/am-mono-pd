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

################################################################################
# MODULE: locker
################################################################################

resource "random_string" "locker_suffix" {
  length  = 8
  special = false
  upper   = false
}

locals {
  locker_id = "${var.project}-${random_string.locker_suffix.result}"
}

module "locker" {
  source = "./modules/locker"

  aws_region          = var.aws_region
  dynamodb_table_name = "locker-${local.locker_id}"
  tags                = var.tags
}

################################################################################
# MODULE: network
################################################################################

resource "random_string" "network_suffix" {
  length  = 8
  special = false
  upper   = false
}

locals {
  network_id = "${var.project}-${random_string.network_suffix.result}"
}

module "network" {
  source = "./modules/network"

  aws_region = var.aws_region
  deployment = var.deployment

  tags = var.tags

  vpc = {
    create_database_internet_gateway_route = {
      production = false
      testing    = true
    }[var.deployment]

    create_database_subnet_route_table = {
      production = false
      testing    = true
    }[var.deployment]

    version = "5.21.0"
  }

  vpc_name = local.network_id
}

################################################################################
# MODULE: tofu_backend
################################################################################

resource "random_string" "mono_tofu_suffix" {
  length  = 12
  special = false
  upper   = false
}

resource "random_string" "terraform_locks_suffix" {
  length  = 12
  special = false
  upper   = false
}

module "tofu_backend" {
  source = "./modules/tofu_backend"

  aws_region                 = var.aws_region
  s3_bucket_name             = "mono-tofu-${random_string.mono_tofu_suffix.result}"
  tags                       = var.tags
  terraform_locks_table_name = "terraform-locks-${random_string.terraform_locks_suffix.result}"
}
