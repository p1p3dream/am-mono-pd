provider "aws" {
  default_tags {
    tags = merge(
      var.tags,
      {
        module = "locker"
      }
    )
  }

  region = var.aws_region
}
