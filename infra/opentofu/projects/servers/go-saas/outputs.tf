output "ecr_repository" {
  value = module.server.ecr_repository
}

output "iam_roles" {
  value = module.server.iam_roles
}

output "load_balancer_target_groups" {
  value = module.server.load_balancer_target_groups
}

output "security_groups" {
  value = module.server.security_groups
}
