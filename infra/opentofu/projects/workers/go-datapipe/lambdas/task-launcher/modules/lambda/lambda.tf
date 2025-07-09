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
      ABODEMINE_DATAPIPE_CONFIG_PATH = var.lambda.config.path
    }
  }

  depends_on = [
    aws_cloudwatch_log_group.main,
  ]
}

resource "aws_lambda_event_source_mapping" "main" {
  event_source_arn = var.sqs.main.arn
  function_name    = aws_lambda_function.main.function_name

  batch_size = var.lambda.config.sqs_event_source.batch_size

  # For FIFO queues, this ensures messages from
  # the same group are processed in order.
  scaling_config {
    # Process maximum_concurrency message groups at a time.
    maximum_concurrency = var.lambda.config.sqs_event_source.maximum_concurrency
  }
}
