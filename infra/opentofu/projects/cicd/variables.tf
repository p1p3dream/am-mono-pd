variable "aws_account_id" {
  description = "AWS account ID."
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

variable "github_organization_name" {
  description = "The name of the GitHub organization."
  type        = string
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
