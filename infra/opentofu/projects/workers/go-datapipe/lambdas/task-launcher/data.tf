data "terraform_remote_state" "main" {
  backend = "s3"

  config = {
    region         = var.aws_region
    bucket         = var.s3_backend_bucket
    key            = var.s3_backend_keys.main
    dynamodb_table = var.s3_backend_table
  }
}

data "terraform_remote_state" "workers_go_datapipe" {
  backend = "s3"

  config = {
    region         = var.aws_region
    bucket         = var.s3_backend_bucket
    key            = var.s3_backend_keys.workers_go_datapipe
    dynamodb_table = var.s3_backend_table
  }
}

data "terraform_remote_state" "workers_shared" {
  backend = "s3"

  config = {
    region         = var.aws_region
    bucket         = var.s3_backend_bucket
    key            = var.s3_backend_keys.workers_shared
    dynamodb_table = var.s3_backend_table
  }
}
