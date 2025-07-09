variable "aws_region" {
  description = "AWS region."
  type        = string
}

variable "dynamodb_table_name" {
  description = "Name of the DynamoDB table."
  type        = string
}

variable "tags" {
  description = "Tags to apply to resources."
  type        = map(string)
}
