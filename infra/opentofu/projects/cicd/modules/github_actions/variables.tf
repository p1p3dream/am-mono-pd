variable "aws_account_id" {
  description = "AWS account ID."
  type        = string
}

variable "aws_region" {
  description = "AWS region."
  type        = string
}

variable "github_org" {
  description = "GitHub organization name"
  type        = string
}

variable "role_name" {
  description = "AWS IAM Role Name"
  type        = string
}

variable "role_policy" {
  description = "AWS IAM Role Policy"
  type        = string
}

variable "tags" {
  description = "Tags to apply to resources."
  type        = map(string)
}
