output "vk" {
  sensitive = true
  value     = aws_elasticache_replication_group.vk
}

output "security_groups" {
  value = {
    elasticache_vk_users_sg = aws_security_group.elasticache_vk_users_sg
  }
}
