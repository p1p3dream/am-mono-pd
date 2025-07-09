variable "aws_region" {
  description = "AWS region."
  type        = string
}

variable "deployment" {
  description = "The deployment target (environment)."
  type        = string
}

variable "ecs_clusters" {
  description = "Map of ECS clusters."
  type = map(object({
    name = string
    setting = object({
      name  = string
      value = string
    })
  }))
}

variable "module_id" {
  description = "The id for this module."
  type        = string
}

variable "tags" {
  description = "Tags to apply to resources."
  type        = map(string)
}
