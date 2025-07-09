variable "aws_account_id" {
  description = "AWS account id."
  type        = string
}

variable "aws_region" {
  description = "AWS region."
  type        = string
}

variable "certs" {
  description = "Certificates."
  type = map(object({
    arn = string
  }))
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

variable "load_balancers" {
  description = "Map of load balancers."
  type = map(object({
    dns_name = string
  }))
}

variable "management_aws_profile" {
  description = "The AWS profile used for the management provider."
  type        = string
}

variable "module_id" {
  description = "The id for this module."
  type        = string
}

variable "module_suffix" {
  description = "Suffix to append to module resources."
  type        = string
}

variable "tags" {
  description = "Tags to apply to resources."
  type        = map(string)
}
