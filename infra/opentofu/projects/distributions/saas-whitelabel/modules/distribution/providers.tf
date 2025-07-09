provider "aws" {
  default_tags {
    tags = merge(
      var.tags,
      {
        module = "distribution"
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
        module = "distribution"
      }
    )
  }

  profile = var.management_aws_profile
  region  = var.aws_region
}

# Provider for Cloudfront certificates.
provider "aws" {
  alias = "us_east_1"

  default_tags {
    tags = merge(
      var.tags,
      {
        module = "distribution"
      }
    )
  }

  region = "us-east-1"
}
