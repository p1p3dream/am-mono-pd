provider "aws" {
  default_tags {
    tags = merge(
      var.tags,
      {
        module = "worker"
      }
    )
  }

  region = var.aws_region
}
