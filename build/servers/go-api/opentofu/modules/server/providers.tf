provider "aws" {
  default_tags {
    tags = merge(
      var.tags,
      {
        module = "server"
      }
    )
  }

  region = var.aws_region
}
