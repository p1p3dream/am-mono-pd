# https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/dynamodb_table.
resource "aws_dynamodb_table" "main" {
  name         = var.dynamodb_table_name
  billing_mode = "PAY_PER_REQUEST"
  hash_key     = "lock_id"

  attribute {
    name = "lock_id"
    type = "S"
  }
}
