output "pg" {
  sensitive = true
  value     = aws_db_instance.pg
}

output "security_groups" {
  value = {
    rds_pg_users_sg = aws_security_group.rds_pg_users_sg
  }
}
