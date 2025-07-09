data "terraform_remote_state" "main" {
  backend = "s3"

  config = {
    region         = var.aws_region
    bucket         = var.s3_backend_bucket
    key            = var.s3_backend_keys.main
    dynamodb_table = var.s3_backend_table
  }
}

data "terraform_remote_state" "workers_go_packer" {
  backend = "s3"

  config = {
    region         = var.aws_region
    bucket         = var.s3_backend_bucket
    key            = var.s3_backend_keys.workers_go_packer
    dynamodb_table = var.s3_backend_table
  }
}
