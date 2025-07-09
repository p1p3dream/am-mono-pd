variable "aws_region" {
  description = "AWS region."
  type        = string
}

variable "deployment" {
  description = "The deployment target (environment)."
  type        = string
}

variable "tags" {
  description = "Tags to apply to resources."
  type        = map(string)
}

variable "vpc" {
  description = "The VPC configuration."
  type = object({
    create_database_internet_gateway_route = bool
    create_database_subnet_route_table     = bool
    version                                = string
  })
}

variable "vpc_name" {
  description = "The name of the VPC."
  type        = string
}
