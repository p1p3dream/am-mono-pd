# https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/apigatewayv2_integration.
resource "aws_apigatewayv2_integration" "main" {
  api_id                 = var.api_gateways.secure_download.id
  integration_type       = "AWS_PROXY"
  integration_uri        = aws_lambda_function.main.invoke_arn
  integration_method     = "POST"
  payload_format_version = "2.0"
}

# https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/apigatewayv2_route.
resource "aws_apigatewayv2_route" "main" {
  api_id    = var.api_gateways.secure_download.id
  route_key = "GET /{proxy+}"
  target    = "integrations/${aws_apigatewayv2_integration.main.id}"
}

# https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/lambda_permission.
resource "aws_lambda_permission" "allow_api_gateway" {
  statement_id  = "AllowAPIGatewayInvoke"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.main.function_name
  principal     = "apigateway.amazonaws.com"
  source_arn    = "${var.api_gateways.secure_download.execution_arn}/*/*"
}
