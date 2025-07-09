resource "aws_security_group" "ssh" {
  name   = "ssh-${var.module_id}"
  vpc_id = var.vpc.id

  ingress {
    description = var.vpc_security_groups.ssh.ingress.description
    from_port   = var.vpc_security_groups.ssh.ingress.from_port
    to_port     = var.vpc_security_groups.ssh.ingress.to_port
    protocol    = var.vpc_security_groups.ssh.ingress.protocol
    cidr_blocks = var.vpc_security_groups.ssh.ingress.cidr_blocks
  }

  egress {
    description = var.vpc_security_groups.ssh.egress.description
    from_port   = var.vpc_security_groups.ssh.egress.from_port
    to_port     = var.vpc_security_groups.ssh.egress.to_port
    protocol    = var.vpc_security_groups.ssh.egress.protocol
    cidr_blocks = var.vpc_security_groups.ssh.egress.cidr_blocks
  }
}
