provider "aws" {
  default_tags {
    tags = merge(
      var.tags,
      {
        module = "build"
      }
    )
  }

  region = var.aws_region
}
