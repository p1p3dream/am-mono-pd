variable "aws_account_id" {
  description = "AWS account id."
  type        = string
}

variable "aws_region" {
  description = "AWS region."
  type        = string
}

variable "deployment" {
  description = "The deployment target (environment)."
  type        = string
}

variable "domains" {
  type = map(object({
    name = string
  }))
}

variable "ecr_repository_names" {
  description = "Name of the ECR repositories used by ECS."
  type        = map(string)
}

variable "ecs_clusters" {
  type = map(object({
    id = string
  }))
}

variable "iam_roles" {
  description = "IAM roles."
  type = map(object({
    name = string
  }))
}

variable "load_balancers" {
  description = "Map of load balancers."
  type = map(object({
    listeners = map(object({
      arn = string
    }))
    security_groups = list(string)
  }))
}

variable "module_id" {
  description = "The id for this module."
  type        = string
}

variable "module_suffix" {
  description = "Suffix to append to module resources."
  type        = string
}

variable "project" {
  description = "The project config."
  type = object({
    containers = map(object({
      name = string
      ports = map(object({
        name = string
        port = number
      }))
    }))
    slug = string
  })
}

variable "tags" {
  description = "Tags to apply to resources."
  type        = map(string)
}

variable "vpc" {
  description = "The VPC configuration."
  type = object({
    id              = string
    private_subnets = list(string)
  })
}
