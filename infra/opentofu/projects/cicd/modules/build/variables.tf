variable "aws_region" {
  description = "AWS region."
  type        = string
}

variable "s3_bucket_names" {
  description = "Names of the S3 buckets used."
  type        = map(string)
}

variable "tags" {
  description = "Tags to apply to resources."
  type        = map(string)
}
