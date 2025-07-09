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

variable "management_aws_profile" {
  description = "The AWS profile used for the management provider."
  type        = string
}

variable "project" {
  description = "The project config."
  type = object({
    slug = string
  })
}

variable "s3_backend_bucket" {
  description = "The S3 bucket used for the Terraform backend."
  type        = string
}

variable "s3_backend_keys" {
  description = "Other backend keys for data fetching."
  type        = map(string)
}

variable "s3_backend_table" {
  description = "The DynamoDB table used for the Terraform backend."
  type        = string
}

variable "tags" {
  description = "A collection of tags for this project, modules, and deployments."
  type        = map(string)
}
