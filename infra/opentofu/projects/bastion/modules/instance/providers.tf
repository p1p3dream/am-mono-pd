provider "aws" {
  default_tags {
    tags = merge(
      var.tags,
      {
        module = "instance"
      }
    )
  }

  region = var.aws_region
}
