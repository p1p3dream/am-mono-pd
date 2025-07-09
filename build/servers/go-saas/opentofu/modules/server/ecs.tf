# https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/ecs_task_definition.
resource "aws_ecs_task_definition" "main" {
  container_definitions    = var.ecs.task_definitions.main.container_definitions
  cpu                      = var.ecs.task_definitions.main.cpu
  execution_role_arn       = var.ecs.task_definitions.main.execution_role_arn
  family                   = var.ecs.task_definitions.main.family
  memory                   = var.ecs.task_definitions.main.memory
  network_mode             = var.ecs.task_definitions.main.network_mode
  requires_compatibilities = var.ecs.task_definitions.main.requires_compatibilities

  runtime_platform {
    operating_system_family = "LINUX"
    cpu_architecture        = "ARM64"
  }
}

# https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/ecs_service.
resource "aws_ecs_service" "main" {
  name    = var.ecs.services.main.name
  cluster = var.ecs.services.main.cluster

  task_definition = aws_ecs_task_definition.main.arn
  desired_count   = var.ecs.services.main.desired_count
  launch_type     = var.ecs.services.main.launch_type

  network_configuration {
    security_groups = var.ecs.services.main.network_configuration.security_groups
    subnets         = var.ecs.services.main.network_configuration.subnets
  }

  load_balancer {
    container_name   = var.ecs.services.main.load_balancer.container_name
    container_port   = var.ecs.services.main.load_balancer.container_port
    target_group_arn = var.ecs.services.main.load_balancer.target_group_arn
  }
}
