output "os_domain" {
  sensitive = true
  value     = module.database.os_domain
}

output "security_groups" {
  value = {
    os_domain_users_sg = module.database.security_groups.os_domain_users_sg
  }
}
