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
    task_definition = object({
      container_definitions    = string
      cpu                      = number
      ephemeral_storage_size   = number
      execution_role_arn       = string
      family                   = string
      memory                   = number
      network_mode             = string
      requires_compatibilities = list(string)
    })
  })
}

variable "tags" {
  description = "Tags to apply to resources."
  type        = map(string)
}
