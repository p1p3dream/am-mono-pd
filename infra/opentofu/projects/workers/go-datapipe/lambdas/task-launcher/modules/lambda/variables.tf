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

variable "dynamodb" {
  description = "DynamoDB config."
  type = map(object({
    allowed_keys = list(string)
    table_arn    = string
  }))
}

variable "ecs_clusters" {
  type = map(object({
    arn = string
    id  = string
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
    allow_pass_role_on = list(string)
    allow_run_task_on  = list(string)
    config = object({
      path = string
      sqs_event_source = object({
        batch_size          = number
        maximum_concurrency = number
      })
    })
    slug = string

    function_name = string
    payload_file  = string
    runtime       = string
    schedule_name = string
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

variable "sqs" {
  description = "SQS config."
  type = map(object({
    arn = string
    url = string
  }))
}

variable "tags" {
  description = "Tags to apply to resources."
  type        = map(string)
}
