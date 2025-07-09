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
  lambda_id = "${var.lambda.slug}-${random_string.lambda_suffix.result}"
}

module "lambda" {
  source = "./modules/lambda"

  aws_account_id = var.aws_account_id
  aws_region     = var.aws_region
  deployment     = var.deployment
  dynamodb       = var.dynamodb

  ecs_clusters = {
    main_fargate = {
      arn = data.terraform_remote_state.workers_shared.outputs.ecs_clusters.main_fargate.arn
      id  = data.terraform_remote_state.workers_shared.outputs.ecs_clusters.main_fargate.id
    }
  }

  iam_roles = {
    eventbridge_lambda_scheduler = {
      name = "lambda-sched-${local.lambda_id}"
    }

    lambda = {
      name = "lambda-${local.lambda_id}"
    }
  }

  lambda = {
    allow_pass_role_on = var.lambda.allow_pass_role_on
    allow_run_task_on  = var.lambda.allow_run_task_on

    config = {
      path = "config.yaml"

      sqs_event_source = {
        batch_size = 4

        # The maximum number of concurrent Lambda function executions
        # that the event source mapping can invoke.
        maximum_concurrency = 2
      }
    }

    slug = var.lambda.slug

    function_name = local.lambda_id
    # This path is relative to the module, not the root of the project.
    payload_file  = "lambda_function_payload.zip"
    runtime       = "provided.al2023"
    schedule_name = "lambda-${local.lambda_id}"
  }

  module_id     = local.lambda_id
  module_suffix = random_string.lambda_suffix.result

  project = var.project

  sqs = {
    main = {
      arn = data.terraform_remote_state.workers_go_datapipe.outputs.sqs_queues.main.arn
      url = data.terraform_remote_state.workers_go_datapipe.outputs.sqs_queues.main.url
    }
  }

  tags = var.tags
}
