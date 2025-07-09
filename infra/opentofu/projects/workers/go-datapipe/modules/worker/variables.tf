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

variable "partner_data" {
  description = "Data about the partners."
  type = map(object({
    s3_bucket_name = string
  }))
}

variable "partner_schedule" {
  description = "Partner task schedules."
  type = map(object({
    schedule_name        = string
    schedule_expression  = string
    sqs_message_body     = string
    sqs_message_group_id = string
  }))
}

variable "project" {
  description = "The project config."
  type = object({
    slug = string
  })
}

variable "sqs_queue_name" {
  description = "The name of the SQS queue."
  type        = string
}

variable "tags" {
  description = "Tags to apply to resources."
  type        = map(string)
}

variable "vpc" {
  description = "The VPC configuration."
  type = object({
    id              = string
    private_subnets = list(string)
  })
}
