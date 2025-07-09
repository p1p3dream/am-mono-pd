output "security_groups" {
  value = {
    aurora_pg_users_sg = aws_security_group.aurora_pg_users_sg
  }
}
