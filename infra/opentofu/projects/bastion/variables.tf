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

variable "project" {
  description = "The project name."
  type        = string
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

variable "user_data_base64" {
  description = "Base64 encoded user data."
  type        = string
}

variable "user_keypair_public_key" {
  description = "The public key for the user keypair."
  type        = string
}
