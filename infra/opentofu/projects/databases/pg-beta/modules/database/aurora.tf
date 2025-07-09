# https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/rds_cluster_parameter_group.
resource "aws_rds_cluster_parameter_group" "aurora_pg_param_group" {
  name   = "aurora-param-group-${var.module_id}"
  family = var.aurora_param_group.family

  # parameter {
  #   name  = "log_statement"
  #   value = "all"
  # }

  parameter {
    # Log statements taking more than 1 seconds.
    name  = "log_min_duration_statement"
    value = "1000"
  }
}

# https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/rds_cluster.
resource "aws_rds_cluster" "aurora_pg_cluster" {
  apply_immediately            = var.aurora_cluster.apply_immediately
  backup_retention_period      = var.aurora_cluster.backup_retention_period
  cluster_identifier           = var.aurora_cluster.cluster_identifier
  database_name                = var.aurora_cluster.database_name
  engine                       = var.aurora_cluster.engine
  engine_version               = var.aurora_cluster.engine_version
  master_username              = var.aurora_cluster.master_username
  performance_insights_enabled = var.aurora_cluster.performance_insights_enabled
  preferred_backup_window      = var.aurora_cluster.preferred_backup_window
  skip_final_snapshot          = var.aurora_cluster.skip_final_snapshot

  availability_zones              = var.vpc.azs
  db_cluster_parameter_group_name = aws_rds_cluster_parameter_group.aurora_pg_param_group.name
  db_subnet_group_name            = var.vpc.database_subnet_group_name
  final_snapshot_identifier       = "final-snapshot-${var.aurora_cluster.cluster_identifier}"
  manage_master_user_password     = true
  storage_encrypted               = true

  vpc_security_group_ids = [
    aws_security_group.aurora_pg_sg.id,
  ]
}

# https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/rds_cluster_instance.
resource "aws_rds_cluster_instance" "aurora_pg_instances" {
  apply_immediately   = var.aurora_cluster.apply_immediately
  count               = var.aurora_instance.count
  instance_class      = var.aurora_instance.instance_class
  publicly_accessible = var.aurora_instance.publicly_accessible

  cluster_identifier           = aws_rds_cluster.aurora_pg_cluster.id
  db_subnet_group_name         = var.vpc.database_subnet_group_name
  engine                       = aws_rds_cluster.aurora_pg_cluster.engine
  engine_version               = aws_rds_cluster.aurora_pg_cluster.engine_version
  identifier                   = "${var.aurora_cluster.cluster_identifier}-${count.index}"
  performance_insights_enabled = var.aurora_cluster.performance_insights_enabled
}
