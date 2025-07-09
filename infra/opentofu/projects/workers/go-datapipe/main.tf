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
# MODULE: worker
################################################################################

resource "random_string" "worker_suffix" {
  length  = 8
  special = false
  upper   = false
}

locals {
  worker_id = "${var.project.slug}-${random_string.worker_suffix.result}"
}

module "worker" {
  source = "./modules/worker"

  aws_account_id = var.aws_account_id
  aws_region     = var.aws_region
  deployment     = var.deployment

  iam_roles = {
    eventbridge_sqs_scheduler = {
      name = "eventbridge-sqs-scheduler-${local.worker_id}"
    }
  }

  module_id     = local.worker_id
  module_suffix = random_string.worker_suffix.result

  partner_data = {
    attom_data = {
      s3_bucket_name = "partner-data-attom-data-${random_string.worker_suffix.result}",
    }

    first_american = {
      s3_bucket_name = "partner-data-first-american-${random_string.worker_suffix.result}",
    }
  }

  partner_schedule = {
    # Run synther only once daily.
    abodemine_synther = {
      schedule_name = "abodemine-synther-${random_string.worker_suffix.result}",
      schedule_expression = {
        production = "cron(0 0 * * ? *)"
        testing    = "cron(30 0 * * ? *)"
      }[var.deployment]
      sqs_message_body = jsonencode({
        partner = "abodemine",
        task    = "synther",
      }),
      sqs_message_group_id = "abodemine"
    }

    attom_data_fetcher = {
      schedule_name = "attom-data-fetcher-${random_string.worker_suffix.result}",
      schedule_expression = {
        production = "cron(0 * * * ? *)"
        testing    = "cron(30 * * * ? *)"
      }[var.deployment]
      sqs_message_body = jsonencode({
        partner = "attom-data",
        task    = "fetcher",
      }),
      sqs_message_group_id = "attom-data"
    }

    attom_data_loader = {
      schedule_name = "attom-data-loader-${random_string.worker_suffix.result}",
      schedule_expression = {
        production = "cron(30 * * * ? *)"
        testing    = "cron(0 * * * ? *)"
      }[var.deployment]
      sqs_message_body = jsonencode({
        partner = "attom-data",
        task    = "loader",
      }),
      sqs_message_group_id = "attom-data"
    }

    first_american_fetcher = {
      schedule_name = "first-american-fetcher-${random_string.worker_suffix.result}",
      schedule_expression = {
        production = "cron(0 * * * ? *)"
        testing    = "cron(30 * * * ? *)"
      }[var.deployment]
      sqs_message_body = jsonencode({
        partner = "first-american",
        task    = "fetcher",
      }),
      sqs_message_group_id = "first-american"
    }

    first_american_loader = {
      schedule_name = "first-american-loader-${random_string.worker_suffix.result}",
      schedule_expression = {
        production = "cron(30 * * * ? *)"
        testing    = "cron(0 * * * ? *)"
      }[var.deployment]
      sqs_message_body = jsonencode({
        partner = "first-american",
        task    = "loader",
      }),
      sqs_message_group_id = "first-american"
    }
  }

  project = var.project

  # SQS queue name will be suffixed with .fifo if set as FIFO.
  sqs_queue_name = local.worker_id
  tags           = var.tags

  vpc = {
    id              = data.terraform_remote_state.main.outputs.vpc.vpc_id
    private_subnets = data.terraform_remote_state.main.outputs.vpc.private_subnets
  }
}
