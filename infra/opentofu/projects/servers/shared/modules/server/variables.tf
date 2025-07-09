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
    name    = string
    zone_id = optional(string, "")
  }))
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

variable "management_aws_profile" {
  description = "The AWS profile used for the management provider."
  type        = string
}

variable "module_id" {
  description = "The id for this module."
  type        = string
}

variable "tags" {
  description = "Tags to apply to resources."
  type        = map(string)
}

variable "vpc" {
  description = "The VPC configuration."
  type = object({
    id             = string
    public_subnets = list(string)
  })
}
