# https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/security_group.

resource "aws_security_group" "rds_pg_users_sg" {
  name   = "rds-pg-users-sg-${var.module_id}"
  vpc_id = var.vpc.id
}

resource "aws_security_group" "rds_pg_sg" {
  name   = "rds-pg-sg-${var.module_id}"
  vpc_id = var.vpc.id

  ingress {
    from_port = 5432
    to_port   = 5432
    protocol  = "tcp"

    security_groups = [
      aws_security_group.rds_pg_users_sg.id,
    ]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}
