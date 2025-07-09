# https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/s3_bucket.
# https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/s3_bucket_lifecycle_configuration.
# https://aws.amazon.com/s3/storage-classes/.

resource "aws_s3_bucket" "secure_download" {
  bucket = var.s3_buckets.secure_download.name
}

resource "aws_s3_bucket_versioning" "secure_download_versioning" {
  bucket = aws_s3_bucket.secure_download.id

  versioning_configuration {
    status = "Enabled"
  }
}
