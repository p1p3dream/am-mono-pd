variable "aws_region" {
  description = "AWS region."
  type        = string
}

variable "deployment" {
  description = "The deployment target (environment)."
  type        = string
}

variable "ecs" {
  description = "The ECS config."
  type = object({
    task_definitions = map(object({
      container_definitions    = string
      cpu                      = number
      execution_role_arn       = string
      family                   = string
      memory                   = number
      network_mode             = string
      requires_compatibilities = list(string)
    }))

    services = map(object({
      name          = string
      cluster       = string
      desired_count = number
      launch_type   = string

      load_balancer = object({
        container_name   = string
        container_port   = number
        target_group_arn = string
      })

      network_configuration = object({
        subnets         = list(string)
        security_groups = list(string)
      })
    }))
  })
}

variable "tags" {
  description = "Tags to apply to resources."
  type        = map(string)
}
