output "api_gateways" {
  description = "API Gateways."
  value       = module.worker.api_gateways
}

output "dynamodb_tables" {
  description = "DynamoDB tables."
  value       = module.worker.dynamodb_tables
}

output "s3_buckets" {
  description = "S3 buckets."
  value       = module.worker.s3_buckets
}
