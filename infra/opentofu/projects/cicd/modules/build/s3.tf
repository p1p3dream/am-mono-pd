# https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/s3_bucket.
# https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/s3_bucket_lifecycle_configuration.
# https://aws.amazon.com/s3/storage-classes/.

resource "aws_s3_bucket" "mono_build" {
  bucket = var.s3_bucket_names.build
}

resource "aws_s3_bucket_lifecycle_configuration" "mono_build_lifecycle" {
  bucket = aws_s3_bucket.mono_build.id

  rule {
    id = "expire"

    filter {
      prefix = "/"
    }

    transition {
      days          = 1
      storage_class = "INTELLIGENT_TIERING"
    }

    expiration {
      days = 30
    }

    status = "Enabled"
  }
}

resource "aws_s3_bucket" "mono_overlay" {
  bucket = var.s3_bucket_names.overlay
}

resource "aws_s3_bucket_versioning" "mono_overlay_versioning" {
  bucket = aws_s3_bucket.mono_overlay.id

  versioning_configuration {
    status = "Enabled"
  }
}
