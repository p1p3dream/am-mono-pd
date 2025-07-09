# https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/lambda_function.
# https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/lambda_event_source_mapping.

resource "aws_lambda_function" "main" {
  function_name = var.lambda.function_name
  role          = aws_iam_role.lambda.arn
  runtime       = var.lambda.runtime
  handler       = "bootstrap"
  architectures = ["arm64"]

  filename         = "${path.module}/${var.lambda.payload_file}"
  source_code_hash = filebase64sha256("${path.module}/${var.lambda.payload_file}")

  environment {
    variables = {
      ABODEMINE_PACKER_CONFIG_PATH = var.lambda.config.path
    }
  }

  depends_on = [
    aws_cloudwatch_log_group.main,
  ]
}
