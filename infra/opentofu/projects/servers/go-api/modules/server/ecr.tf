# https://registry.terraform.io/providers/hashicorp/aws/latest/docs/data-sources/ecr_repository.
resource "aws_ecr_repository" "main" {
  for_each = var.ecr_repository_names

  name                 = each.value
  image_tag_mutability = "MUTABLE"

  image_scanning_configuration {
    scan_on_push = true
  }
}
