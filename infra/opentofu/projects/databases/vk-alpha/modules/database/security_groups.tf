# https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/security_group.

resource "aws_security_group" "elasticache_vk_users_sg" {
  name   = "elasticache-vk-users-sg-${var.module_id}"
  vpc_id = var.vpc.id
}

resource "aws_security_group" "elasticache_vk_sg" {
  name   = "elasticache-vk-sg-${var.module_id}"
  vpc_id = var.vpc.id

  ingress {
    from_port = 6379
    to_port   = 6379
    protocol  = "tcp"

    security_groups = [
      aws_security_group.elasticache_vk_users_sg.id,
    ]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}
