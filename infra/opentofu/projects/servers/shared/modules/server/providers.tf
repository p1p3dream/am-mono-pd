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

provider "aws" {
  alias = "management"

  default_tags {
    tags = merge(
      var.tags,
      {
        module = "server"
      }
    )
  }

  profile = var.management_aws_profile
  region  = var.aws_region
}
