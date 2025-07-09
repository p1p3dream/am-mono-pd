data "terraform_remote_state" "main" {
  backend = "s3"

  config = {
    region         = var.aws_region
    bucket         = var.s3_backend_bucket
    key            = var.s3_backend_keys.main
    dynamodb_table = var.s3_backend_table
  }
}

data "terraform_remote_state" "servers_shared" {
  backend = "s3"

  config = {
    region         = var.aws_region
    bucket         = var.s3_backend_bucket
    key            = var.s3_backend_keys.servers_shared
    dynamodb_table = var.s3_backend_table
  }
}
