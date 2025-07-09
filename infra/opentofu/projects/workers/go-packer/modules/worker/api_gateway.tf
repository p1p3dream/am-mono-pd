# https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/apigatewayv2_api.
resource "aws_apigatewayv2_api" "secure_download" {
  name          = "secure-download-${var.module_suffix}"
  protocol_type = "HTTP"
}

# https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/apigatewayv2_stage.
resource "aws_apigatewayv2_stage" "secure_download" {
  api_id      = aws_apigatewayv2_api.secure_download.id
  name        = "v1"
  auto_deploy = true
}
