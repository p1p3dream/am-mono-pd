variable "aws_region" {
  description = "AWS region."
  type        = string
}

variable "s3_bucket_name" {
  description = "Name of the S3 bucket used for Terraform state storage."
  type        = string
}

variable "tags" {
  description = "Tags to apply to resources."
  type        = map(string)
}

variable "terraform_locks_table_name" {
  description = "Name of the DynamoDB table used for Terraform state locking."
  type        = string
}
