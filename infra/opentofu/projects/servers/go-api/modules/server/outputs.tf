output "ecr_repository" {
  value = {
    main = aws_ecr_repository.main
  }
}

output "iam_roles" {
  value = {
    main = aws_iam_role.main
  }
}

output "load_balancer_target_groups" {
  value = {
    main = aws_lb_target_group.main
  }
}

output "security_groups" {
  value = {
    main = aws_security_group.main
  }
}
