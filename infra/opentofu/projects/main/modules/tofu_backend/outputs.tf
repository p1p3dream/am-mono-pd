output "dynamodb_table" {
  value = aws_dynamodb_table.terraform_locks
}

output "s3_bucket" {
  value = aws_s3_bucket.mono_tofu
}
