# https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/dynamodb_table.
resource "aws_dynamodb_table" "secure_download" {
  name         = var.dynamodb_tables.secure_download.name
  billing_mode = "PAY_PER_REQUEST"
  hash_key     = "token"

  attribute {
    name = "token"
    type = "S"
  }
}
