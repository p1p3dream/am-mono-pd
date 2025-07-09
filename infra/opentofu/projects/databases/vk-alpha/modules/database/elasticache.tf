# https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/elasticache_replication_group.
resource "aws_elasticache_replication_group" "vk" {
  apply_immediately        = var.elasticache_replication_group.apply_immediately
  description              = var.elasticache_replication_group.description
  engine                   = var.elasticache_replication_group.engine
  engine_version           = var.elasticache_replication_group.engine_version
  maintenance_window       = var.elasticache_replication_group.maintenance_window
  node_type                = var.elasticache_replication_group.node_type
  num_node_groups          = var.elasticache_replication_group.num_node_groups
  parameter_group_name     = var.elasticache_replication_group.parameter_group_name
  port                     = var.elasticache_replication_group.port
  replication_group_id     = var.elasticache_replication_group.replication_group_id
  snapshot_retention_limit = var.elasticache_replication_group.snapshot_retention_limit
  snapshot_window          = var.elasticache_replication_group.snapshot_window

  security_group_ids = [
    aws_security_group.elasticache_vk_sg.id,
  ]

  subnet_group_name = var.vpc.elasticache_subnet_group_name
}
