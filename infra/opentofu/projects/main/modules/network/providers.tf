provider "aws" {
  default_tags {
    tags = merge(
      var.tags,
      {
        module = "network"
      }
    )
  }

  region = var.aws_region
}
