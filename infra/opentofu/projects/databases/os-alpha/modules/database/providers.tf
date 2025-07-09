provider "aws" {
  default_tags {
    tags = merge(
      var.tags,
      {
        module = "database"
      }
    )
  }

  region = var.aws_region
}
