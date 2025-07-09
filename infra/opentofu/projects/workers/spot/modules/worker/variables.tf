variable "aws_account_id" {
  description = "AWS account id."
  type        = string
}

variable "aws_region" {
  description = "AWS region."
  type        = string
}

variable "config" {
  description = "Local config for the worker."
  type = object({
    vars = object({
      block_devices = object({
        root = object({
          delete_on_termination = string
          volume_type           = string
          volume_size           = number
          iops                  = number
          throughput            = number
        })

        persistent = object({
          delete_on_termination = string
          device_name           = string
          mount_point           = string
          volume_type           = string
          volume_size           = number
          iops                  = number
          throughput            = number
        })
      })

      debian_version        = string
      ebs_availability_zone = string
      ec2_instance_type     = string
      spot_instance         = bool
      ssh_port              = number
    })
  })
}

variable "deployment" {
  description = "The deployment target (environment)."
  type        = string
}

variable "module_id" {
  description = "The id for this module."
  type        = string
}

variable "subnet_id" {
  description = "The subnet id to deploy the worker to."
  type        = string
}

variable "tags" {
  description = "Tags to apply to resources."
  type        = map(string)
}

variable "user_data_base64" {
  description = "Base64 encoded user data."
  type        = string
}

variable "user_keypair_key_name" {
  description = "The key name for the user keypair."
  type        = string
}

variable "user_keypair_public_key" {
  description = "The public key for the user keypair."
  type        = string
}

variable "vpc" {
  description = "The VPC configuration."
  type = object({
    id = string
  })
}

variable "vpc_security_groups" {
  description = "Security groups to apply to resources."

  type = map(object({
    ingress = object({
      description = string
      from_port   = number
      to_port     = number
      protocol    = string
      cidr_blocks = list(string)
    })
    egress = object({
      description = string
      from_port   = number
      to_port     = number
      protocol    = string
      cidr_blocks = list(string)
    })
  }))
}

variable "vpc_security_group_ids" {
  description = "The security group ids to apply to the worker."
  type        = list(string)
}
