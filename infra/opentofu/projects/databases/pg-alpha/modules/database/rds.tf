# https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/db_instance.
resource "aws_db_instance" "pg" {
  allocated_storage            = var.rds.allocated_storage
  apply_immediately            = var.rds.apply_immediately
  backup_retention_period      = var.rds.backup_retention_period
  backup_window                = var.rds.backup_window
  db_name                      = var.rds.db_name
  engine                       = var.rds.engine
  engine_version               = var.rds.engine_version
  identifier                   = var.rds.identifier
  instance_class               = var.rds.instance_class
  iops                         = var.rds.iops
  maintenance_window           = var.rds.maintenance_window
  manage_master_user_password  = var.rds.manage_master_user_password
  max_allocated_storage        = var.rds.max_allocated_storage
  monitoring_interval          = var.rds.monitoring_interval
  multi_az                     = var.rds.multi_az
  performance_insights_enabled = var.rds.performance_insights_enabled
  publicly_accessible          = var.rds.publicly_accessible
  skip_final_snapshot          = var.rds.skip_final_snapshot
  storage_encrypted            = var.rds.storage_encrypted
  storage_throughput           = var.rds.storage_throughput
  storage_type                 = var.rds.storage_type
  username                     = var.rds.username

  db_subnet_group_name = var.vpc.database_subnet_group_name
  monitoring_role_arn  = aws_iam_role.rds_pg_role.arn

  vpc_security_group_ids = [
    aws_security_group.rds_pg_sg.id,
  ]
}
