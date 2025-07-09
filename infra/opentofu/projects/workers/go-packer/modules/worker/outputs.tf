output "api_gateways" {
  description = "API Gateways."
  value = {
    secure_download = aws_apigatewayv2_api.secure_download
  }
}

output "dynamodb_tables" {
  description = "DynamoDB tables."
  value = {
    secure_download = aws_dynamodb_table.secure_download
  }
}

output "s3_buckets" {
  description = "S3 buckets."
  value = {
    secure_download = aws_s3_bucket.secure_download
  }
}
