output "ecr_repositories" {
  description = "Map of ECR repositories."
  value = {
    for k, v in var.ecr_repository_names : k => aws_ecr_repository.main[k]
  }
}

output "iam_roles" {
  description = "Map of IAM roles."
  value = {
    main = aws_iam_role.main
  }
}

output "security_groups" {
  description = "Map of security groups."
  value = {
    main = aws_security_group.main
  }
}
