# https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/s3_bucket.
# https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/s3_bucket_lifecycle_configuration.
# https://aws.amazon.com/s3/storage-classes/.

resource "aws_s3_bucket" "mono_tofu" {
  bucket = var.s3_bucket_name
}

resource "aws_s3_bucket_versioning" "mono_tofu_versioning" {
  bucket = aws_s3_bucket.mono_tofu.id

  versioning_configuration {
    status = "Enabled"
  }
}
