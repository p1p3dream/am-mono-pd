output "security_groups" {
  value = {
    aurora_pg_users_sg = module.database.security_groups.aurora_pg_users_sg
  }
}
