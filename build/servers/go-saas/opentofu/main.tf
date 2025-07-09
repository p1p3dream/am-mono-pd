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
# MODULE: server
################################################################################

resource "random_string" "server_suffix" {
  length  = 8
  special = false
  upper   = false
}

locals {
  server_id = "${var.project.slug}-${random_string.server_suffix.result}"
}

module "server" {
  source = "./modules/server"

  aws_region = var.aws_region
  deployment = var.deployment

  ecs = {
    task_definitions = {
      main = {
        family                   = local.server_id
        cpu                      = var.project.containers.main.cpu
        execution_role_arn       = var.iam_roles.main.arn
        memory                   = var.project.containers.main.memory
        network_mode             = "awsvpc"
        requires_compatibilities = ["FARGATE"]

        container_definitions = jsonencode([
          {
            name        = var.project.containers.main.name
            image       = var.project.containers.main.image
            cpu         = var.project.containers.main.cpu
            memory      = var.project.containers.main.memory
            entrypoint  = ["/app/bin/server"]
            command     = ["listen"]
            essential   = true
            environment = var.project.containers.main.env
            portMappings = [
              {
                containerPort = var.project.containers.main.ports.http.port
                hostPort      = var.project.containers.main.ports.http.port
              }
            ]
            logConfiguration = {
              logDriver = "awslogs"
              options = {
                "awslogs-group"         = "/ecs/${local.server_id}"
                "awslogs-region"        = var.aws_region
                "awslogs-stream-prefix" = "ecs"
                "awslogs-create-group"  = "true"
              }
            }
          }
        ])
      }
    }

    services = {
      main = {
        name    = local.server_id
        cluster = var.ecs.clusters.main_fargate.arn

        desired_count = {
          production = 3
          testing    = 2
        }[var.deployment]

        launch_type           = "FARGATE"
        load_balancer         = var.ecs.services.main.load_balancer
        network_configuration = var.ecs.services.main.network_configuration
      }
    }
  }

  tags = var.tags
}
