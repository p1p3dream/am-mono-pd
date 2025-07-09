output "vk" {
  sensitive = true
  value     = module.database.vk
}

output "security_groups" {
  value = {
    elasticache_vk_users_sg = module.database.security_groups.elasticache_vk_users_sg
  }
}
