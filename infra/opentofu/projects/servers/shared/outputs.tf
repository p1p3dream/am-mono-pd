output "certs" {
  sensitive = true
  value     = module.server.certs
}

output "ecs_clusters" {
  value = module.server.ecs_clusters
}

output "load_balancers" {
  value = module.server.load_balancers
}
