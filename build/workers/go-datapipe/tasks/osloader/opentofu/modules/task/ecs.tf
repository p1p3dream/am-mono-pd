# https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/ecs_task_definition.
resource "aws_ecs_task_definition" "main" {
  container_definitions    = var.ecs.task_definition.container_definitions
  cpu                      = var.ecs.task_definition.cpu
  execution_role_arn       = var.ecs.task_definition.execution_role_arn
  family                   = var.ecs.task_definition.family
  memory                   = var.ecs.task_definition.memory
  network_mode             = var.ecs.task_definition.network_mode
  requires_compatibilities = var.ecs.task_definition.requires_compatibilities
  task_role_arn            = var.ecs.task_definition.execution_role_arn

  dynamic "ephemeral_storage" {
    for_each = var.ecs.task_definition.ephemeral_storage_size != null ? [var.ecs.task_definition.ephemeral_storage_size] : []
    content {
      size_in_gib = ephemeral_storage.value
    }
  }

  runtime_platform {
    operating_system_family = "LINUX"
    cpu_architecture        = "ARM64"
  }
}
