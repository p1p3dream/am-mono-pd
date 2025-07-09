# https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/security_group.html.
resource "aws_security_group" "main_lb" {
  name        = "main-lb-${var.module_id}"
  description = "Security group for the main LB."
  vpc_id      = var.vpc.id

  ingress {
    from_port   = 443
    to_port     = 443
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

# https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/lb.
resource "aws_lb" "main" {
  name               = "main-lb-${var.module_id}"
  internal           = false
  load_balancer_type = "application"


  security_groups = [aws_security_group.main_lb.id]
  subnets         = var.vpc.public_subnets

  enable_deletion_protection = true
}

resource "aws_lb_listener" "main_http" {
  load_balancer_arn = aws_lb.main.arn
  port              = "80"
  protocol          = "HTTP"

  default_action {
    type = "redirect"
    redirect {
      port        = "443"
      protocol    = "HTTPS"
      status_code = "HTTP_301"
    }
  }
}

# https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/lb_listener.
resource "aws_lb_listener" "main_https" {
  load_balancer_arn = aws_lb.main.arn
  port              = "443"
  protocol          = "HTTPS"
  ssl_policy        = "ELBSecurityPolicy-2016-08"

  # Default certificate.
  certificate_arn = aws_acm_certificate.abodemine_main.arn

  default_action {
    type = "fixed-response"
    fixed_response {
      content_type = "text/plain"
      message_body = "No Abode to Mine here."
      status_code  = "404"
    }
  }
}
