output "ecr_repositories" {
  value = module.task.ecr_repositories
}

output "iam_roles" {
  value = module.task.iam_roles
}

output "security_groups" {
  value = module.task.security_groups
}
