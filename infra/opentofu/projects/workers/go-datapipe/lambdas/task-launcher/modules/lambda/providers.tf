provider "aws" {
  default_tags {
    tags = merge(
      var.tags,
      {
        module = "lambda"
      }
    )
  }

  region = var.aws_region
}
