# https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/cloudwatch_log_group.

resource "aws_cloudwatch_log_group" "main" {
  name              = "/aws/lambda/${var.lambda.function_name}"
  retention_in_days = 30
}
