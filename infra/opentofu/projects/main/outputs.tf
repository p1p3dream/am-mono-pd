output "dynamodb_table_locker" {
  value = module.locker.dynamodb_table_main
}

output "vpc" {
  value = module.network.vpc
}
