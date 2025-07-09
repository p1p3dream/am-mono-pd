# https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/security_group.

resource "aws_security_group" "aurora_pg_users_sg" {
  name   = "aurora-pg-users-sg-${var.module_id}"
  vpc_id = var.vpc.id
}

resource "aws_security_group" "aurora_pg_ext_users_sg" {
  name   = "aurora-pg-ext-users-sg-${var.module_id}"
  vpc_id = var.vpc.id

  # ingress {
  #   from_port = 5432
  #   to_port   = 5432
  #   protocol  = "tcp"

  #   cidr_blocks = ["0.0.0.0/0"]
  # }

  # egress {
  #   from_port   = 0
  #   to_port     = 0
  #   protocol    = "tcp"
  #   cidr_blocks = ["0.0.0.0/0"]
  # }
}

resource "aws_security_group" "aurora_pg_sg" {
  name   = "aurora-pg-sg-${var.module_id}"
  vpc_id = var.vpc.id

  ingress {
    from_port = 5432
    to_port   = 5432
    protocol  = "tcp"

    cidr_blocks = {
      production = concat(
        [
          # MotherDuck IPs.
          # "44.209.39.178/32",
          # "52.22.9.251/32"
        ]
      )
      testing = concat(
        ["0.0.0.0/0"]
      ),
    }[var.deployment]

    security_groups = [
      aws_security_group.aurora_pg_users_sg.id,
      aws_security_group.aurora_pg_ext_users_sg.id,
    ]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }
}
