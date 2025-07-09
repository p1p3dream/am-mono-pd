provider "aws" {
  default_tags {
    tags = merge(
      var.tags,
      {
        module = "task"
      }
    )
  }

  region = var.aws_region
}
