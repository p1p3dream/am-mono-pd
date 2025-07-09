provider "aws" {
  default_tags {
    tags = merge(
      var.tags,
      {
        module = "tofu_backend"
      }
    )
  }

  region = var.aws_region
}
