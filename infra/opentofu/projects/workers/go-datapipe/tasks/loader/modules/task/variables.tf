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

variable "dynamodb_tables" {
  description = "DynamoDB tables."
  type = map(object({
    allowed_keys = list(string)
    table_arn    = string
  }))
}

variable "ecr_repository_names" {
  description = "Name of the ECR repositories used by ECS."
  type        = map(string)
}

variable "iam_roles" {
  description = "IAM roles."
  type = map(object({
    name = string
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

variable "s3_buckets" {
  description = "S3 buckets."
  type        = map(list(string))
}

variable "security_groups" {
  description = "Security groups."
  type = map(object({
    name = string
  }))
}

variable "tags" {
  description = "Tags to apply to resources."
  type        = map(string)
}

variable "task" {
  description = "The task config."
  type = object({
    slug = string
  })
}

variable "vpc" {
  description = "The VPC configuration."
  type = object({
    id              = string
    private_subnets = list(string)
  })
}
