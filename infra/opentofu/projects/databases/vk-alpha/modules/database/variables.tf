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

variable "elasticache_replication_group" {
  description = "Elasticache config."
  type = object({
    apply_immediately        = bool
    description              = string
    engine                   = string
    engine_version           = string
    maintenance_window       = string
    node_type                = string
    num_node_groups          = number
    parameter_group_name     = string
    port                     = number
    replicas_per_node_group  = number
    replication_group_id     = string
    snapshot_retention_limit = optional(number)
    snapshot_window          = optional(string)
  })
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
    elasticache_subnet_group_name = string
    id                            = string
  })
}
