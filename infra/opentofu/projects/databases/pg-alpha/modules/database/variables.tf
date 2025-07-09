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

variable "module_id" {
  description = "The id for this module."
  type        = string
}

variable "rds" {
  description = "RDS config."
  type = object({
    allocated_storage            = number
    apply_immediately            = bool
    backup_retention_period      = number
    backup_window                = string
    db_name                      = string
    engine                       = string
    engine_version               = string
    identifier                   = string
    instance_class               = string
    iops                         = optional(number)
    maintenance_window           = string
    manage_master_user_password  = bool
    max_allocated_storage        = number
    monitoring_interval          = number
    multi_az                     = bool
    performance_insights_enabled = bool
    publicly_accessible          = bool
    skip_final_snapshot          = bool
    storage_encrypted            = bool
    storage_throughput           = optional(number)
    storage_type                 = string
    username                     = string
  })
}

variable "tags" {
  description = "Tags to apply to resources."
  type        = map(string)
}

variable "vpc" {
  description = "The VPC configuration."
  type = object({
    database_subnet_group_name = string
    id                         = string
  })
}
