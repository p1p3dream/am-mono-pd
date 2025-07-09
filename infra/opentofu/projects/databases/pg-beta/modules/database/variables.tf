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

variable "aurora_cluster" {
  description = "Aurora cluster configuration."
  type = object({
    apply_immediately            = bool
    backup_retention_period      = number
    cluster_identifier           = string
    database_name                = string
    engine                       = string
    engine_version               = string
    master_username              = string
    performance_insights_enabled = bool
    preferred_backup_window      = string
    skip_final_snapshot          = bool
  })
}

variable "aurora_instance" {
  description = "Aurora instance configuration."
  type = object({
    apply_immediately   = bool
    count               = number
    instance_class      = string
    publicly_accessible = bool
  })
}

variable "aurora_param_group" {
  description = "Aurora parameter group configuration."
  type = object({
    family = string
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
    azs                          = list(string)
    database_subnets_cidr_blocks = list(string)
    database_subnet_group_name   = string
    id                           = string
    private_subnets_cidr_blocks  = list(string)
  })
}
