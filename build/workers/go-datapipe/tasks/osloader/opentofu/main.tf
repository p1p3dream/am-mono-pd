terraform {
  backend "s3" {
    region         = var.aws_region
    bucket         = var.s3_backend_bucket
    key            = var.s3_backend_keys.task
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
# MODULE: task
################################################################################

resource "random_string" "task_suffix" {
  length  = 8
  special = false
  upper   = false
}

locals {
  task_id = "${var.task.slug}-${random_string.task_suffix.result}"
}

module "task" {
  source = "./modules/task"

  aws_region = var.aws_region
  deployment = var.deployment

  ecs = {
    task_definition = {
      family                   = local.task_id
      cpu                      = var.task.containers.main.cpu
      ephemeral_storage_size   = var.task.ephemeral_storage_size
      execution_role_arn       = var.iam_roles.main.arn
      memory                   = var.task.containers.main.memory
      network_mode             = "awsvpc"
      requires_compatibilities = ["FARGATE"]

      container_definitions = jsonencode([
        {
          name        = var.task.containers.main.name
          image       = var.task.containers.main.image
          cpu         = var.task.containers.main.cpu
          memory      = var.task.containers.main.memory
          entrypoint  = ["/app/bin/worker"]
          command     = []
          essential   = true
          environment = var.task.containers.main.env
          logConfiguration = {
            logDriver = "awslogs"
            options = {
              "awslogs-group"         = "/ecs/${local.task_id}"
              "awslogs-region"        = var.aws_region
              "awslogs-stream-prefix" = "ecs"
              "awslogs-create-group"  = "true"
            }
          }
        }
      ])
    }
  }

  tags = var.tags
}
