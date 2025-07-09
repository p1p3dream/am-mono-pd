# https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/s3_bucket.
# https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/s3_bucket_lifecycle_configuration.
# https://aws.amazon.com/s3/storage-classes/.

resource "aws_s3_bucket" "mono_cache" {
  bucket = "mono-cache-${var.module_suffix}"
}

resource "aws_s3_bucket_versioning" "mono_cache_versioning" {
  bucket = aws_s3_bucket.mono_cache.id

  versioning_configuration {
    status = "Enabled"
  }
}
