provider "aws" {
  default_tags {
    tags = merge(
      var.tags,
      {
        module = "bucket"
      }
    )
  }

  region = var.aws_region
}
