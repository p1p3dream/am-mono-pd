variable "api_gateways" {
  description = "API Gateways."
  type = map(object({
    execution_arn = string
    id            = string
  }))
}

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
    arn = string
  }))
}

variable "iam_roles" {
  description = "IAM roles."
  type = map(object({
    name = string
  }))
}

variable "lambda" {
  description = "The lambda config."
  type = object({
    config = object({
      path = string
    })
    slug = string

    function_name = string
    payload_file  = string
    runtime       = string
  })
}

variable "module_id" {
  description = "The id for this module."
  type        = string
}

variable "module_suffix" {
  description = "Suffix to append to module resources."
  type        = string
}

variable "project" {
  description = "The project config."
  type = object({
    slug = string
  })
}

variable "s3_buckets" {
  description = "S3 buckets."
  type = map(object({
    arn = string
  }))
}

variable "tags" {
  description = "Tags to apply to resources."
  type        = map(string)
}
