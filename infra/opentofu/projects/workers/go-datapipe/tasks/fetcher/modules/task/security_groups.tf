# https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/security_group.html.

resource "aws_security_group" "main" {
  name   = var.security_groups.main.name
  vpc_id = var.vpc.id

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}
