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
    name = string
  }))
}

variable "main_s3_backend_key" {
  description = "The S3 key used for the 'main' Terraform backend."
  type        = string
}

variable "project" {
  description = "The project config."
  type = object({
    containers = map(object({
      name = string
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

variable "servers_shared_s3_backend_key" {
  description = "The S3 key used for the 'servers-shared' Terraform backend."
  type        = string
}

variable "tags" {
  description = "A collection of tags for this project, modules, and deployments."
  type        = map(string)
}
