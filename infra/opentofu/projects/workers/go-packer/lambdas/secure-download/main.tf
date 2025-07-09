terraform {
  backend "s3" {
    region         = var.aws_region
    bucket         = var.s3_backend_bucket
    key            = var.s3_backend_keys.lambda
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
# MODULE: lambda
################################################################################

resource "random_string" "lambda_suffix" {
  length  = 8
  special = false
  upper   = false
}

locals {
  lambda_id         = "${var.lambda.slug}-${random_string.lambda_suffix.result}"
  workers_go_packer = data.terraform_remote_state.workers_go_packer.outputs
}

module "lambda" {
  source = "./modules/lambda"

  api_gateways = {
    secure_download = {
      execution_arn = local.workers_go_packer.api_gateways.secure_download.execution_arn
      id            = local.workers_go_packer.api_gateways.secure_download.id
    }
  }

  aws_account_id = var.aws_account_id
  aws_region     = var.aws_region
  deployment     = var.deployment

  dynamodb_tables = {
    secure_download = {
      arn = local.workers_go_packer.dynamodb_tables.secure_download.arn
    }
  }

  iam_roles = {
    lambda = {
      name = "lambda-${local.lambda_id}"
    }
  }

  lambda = {
    config = {
      path = "config.yaml"
    }

    slug = var.lambda.slug

    function_name = local.lambda_id
    # This path is relative to the module, not the root of the project.
    payload_file = "lambda_function_payload.zip"
    runtime      = "provided.al2023"
  }

  module_id     = local.lambda_id
  module_suffix = random_string.lambda_suffix.result

  project = var.project

  s3_buckets = {
    secure_download = {
      arn = local.workers_go_packer.s3_buckets.secure_download.arn
    }
  }

  tags = var.tags
}
