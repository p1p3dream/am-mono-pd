output "os_domain" {
  sensitive = true
  value     = aws_opensearch_domain.os
}

output "security_groups" {
  value = {
    os_domain_users_sg = aws_security_group.os_domain_users_sg
  }
}
