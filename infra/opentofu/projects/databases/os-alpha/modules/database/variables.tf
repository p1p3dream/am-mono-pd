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

variable "os_domain" {
  description = "OpenSearch domain configuration."
  type = object({
    domain_name    = string
    engine_version = string

    cluster_config = object({
      instance_type  = string
      instance_count = number

      dedicated_master_enabled = bool
      dedicated_master_type    = string
      dedicated_master_count   = number

      warm_enabled           = bool
      zone_awareness_enabled = bool

      zone_awareness_config = optional(object({
        availability_zone_count = number
      }))
    })

    ebs_options = object({
      ebs_enabled = bool
      volume_size = number
      volume_type = string
      iops        = number
      throughput  = number
    })

    encrypt_at_rest = object({
      enabled = bool
    })

    node_to_node_encryption = object({
      enabled = bool
    })

    domain_endpoint_options = optional(object({
      enforce_https       = bool
      tls_security_policy = string
    }))

    advanced_security_options = optional(object({
      enabled                        = bool
      internal_user_database_enabled = bool
      master_user_options = object({
        master_user_name     = string
        master_user_password = string
      })
    }))
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
    database_subnets = list(string)
    id               = string
  })
}
