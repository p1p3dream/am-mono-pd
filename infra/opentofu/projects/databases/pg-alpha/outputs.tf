output "pg" {
  sensitive = true
  value     = module.database.pg
}

output "security_groups" {
  value = {
    rds_pg_users_sg = module.database.security_groups.rds_pg_users_sg
  }
}
