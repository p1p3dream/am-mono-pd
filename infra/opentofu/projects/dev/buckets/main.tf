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
# MODULE: bucket
################################################################################

resource "random_string" "bucket_suffix" {
  length  = 8
  special = false
  upper   = false
}

locals {
  bucket_id = "${var.project.slug}-${random_string.bucket_suffix.result}"
}

module "bucket" {
  source = "./modules/bucket"

  aws_account_id = var.aws_account_id
  aws_region     = var.aws_region
  deployment     = var.deployment
  module_id      = local.bucket_id
  module_suffix  = random_string.bucket_suffix.result
  tags           = var.tags
}
