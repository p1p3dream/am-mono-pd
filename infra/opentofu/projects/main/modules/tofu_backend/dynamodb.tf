# https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/dynamodb_table.
resource "aws_dynamodb_table" "terraform_locks" {
  name         = var.terraform_locks_table_name
  billing_mode = "PAY_PER_REQUEST"
  hash_key     = "LockID"

  attribute {
    name = "LockID"
    type = "S"
  }
}
