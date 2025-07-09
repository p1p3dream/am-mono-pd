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

variable "ecs" {
  description = "The ECS config."
  type = object({
    clusters = map(object({
      arn = string
    }))

    services = map(object({
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

variable "iam_roles" {
  description = "IAM roles."
  type = map(object({
    arn = string
  }))
}

variable "project" {
  description = "The project config."
  type = object({
    containers = map(object({
      cpu = number
      env = list(object({
        name  = string
        value = string
      }))
      memory = number
      name   = string
      image  = string
      ports = map(object({
        name = string
        port = number
      }))
    }))
    slug = string
  })
}

variable "s3_backend_bucket" {
  description = "The S3 bucket used for the Terraform backend."
  type        = string
}

variable "s3_backend_key" {
  description = "The S3 key used for the Terraform backend."
  type        = string
}

variable "s3_backend_table" {
  description = "The DynamoDB table used for the Terraform backend."
  type        = string
}

variable "tags" {
  description = "A collection of tags for this project, modules, and deployments."
  type        = map(string)
}
