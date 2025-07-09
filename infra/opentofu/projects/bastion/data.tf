data "terraform_remote_state" "databases_os_alpha" {
  backend = "s3"

  config = {
    region         = var.aws_region
    bucket         = var.s3_backend_bucket
    key            = var.s3_backend_keys.databases_os_alpha
    dynamodb_table = var.s3_backend_table
  }
}

data "terraform_remote_state" "databases_pg_alpha" {
  backend = "s3"

  config = {
    region         = var.aws_region
    bucket         = var.s3_backend_bucket
    key            = var.s3_backend_keys.databases_pg_alpha
    dynamodb_table = var.s3_backend_table
  }
}

data "terraform_remote_state" "databases_pg_beta" {
  backend = "s3"

  config = {
    region         = var.aws_region
    bucket         = var.s3_backend_bucket
    key            = var.s3_backend_keys.databases_pg_beta
    dynamodb_table = var.s3_backend_table
  }
}

data "terraform_remote_state" "databases_vk_alpha" {
  backend = "s3"

  config = {
    region         = var.aws_region
    bucket         = var.s3_backend_bucket
    key            = var.s3_backend_keys.databases_vk_alpha
    dynamodb_table = var.s3_backend_table
  }
}

data "terraform_remote_state" "main" {
  backend = "s3"

  config = {
    region         = var.aws_region
    bucket         = var.s3_backend_bucket
    key            = var.s3_backend_keys.main
    dynamodb_table = var.s3_backend_table
  }
}
