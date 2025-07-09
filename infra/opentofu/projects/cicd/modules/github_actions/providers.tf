provider "aws" {
  default_tags {
    tags = merge(
      var.tags,
      {
        module = "github_actions"
      }
    )
  }

  region = var.aws_region
}
