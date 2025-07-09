# https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/s3_bucket.
# https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/s3_bucket_lifecycle_configuration.
# https://aws.amazon.com/s3/storage-classes/.

resource "aws_s3_bucket" "partner_data" {
  for_each = var.partner_data

  bucket = each.value.s3_bucket_name
}

resource "aws_s3_bucket_versioning" "partner_data_versioning" {
  for_each = var.partner_data

  bucket = aws_s3_bucket.partner_data[each.key].id

  versioning_configuration {
    status = "Enabled"
  }
}

resource "aws_s3_bucket_lifecycle_configuration" "partner_data_lifecycle" {
  for_each = var.partner_data

  bucket = aws_s3_bucket.partner_data[each.key].id

  rule {
    id     = "archive"
    status = "Enabled"

    filter {
      prefix = "/"
    }

    transition {
      days          = 180
      storage_class = "STANDARD_IA"
    }

    transition {
      days          = 360
      storage_class = "GLACIER_IR"
    }

    transition {
      days          = 720
      storage_class = "DEEP_ARCHIVE"
    }
  }
}
