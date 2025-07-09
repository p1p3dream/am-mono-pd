# https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/lb_target_group.
resource "aws_lb_target_group" "main" {
  name             = "main-${var.module_id}"
  port             = 443
  protocol         = "HTTPS"
  protocol_version = "HTTP1"
  vpc_id           = var.vpc.id
  target_type      = "ip"

  health_check {
    protocol            = "HTTPS"
    port                = var.project.containers.main.ports.http.port
    path                = "/api/685f26df-dc42-4009-8784-63bd162d0255"
    healthy_threshold   = 2
    unhealthy_threshold = 10
  }
}

# https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/lb_listener_rule.
resource "aws_lb_listener_rule" "main" {
  listener_arn = var.load_balancers.main.listeners.https.arn
  priority     = 200

  action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.main.arn
  }

  condition {
    host_header {
      values = [var.domains.abodemine_saas.name]
    }
  }
}

# https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/security_group.html.
resource "aws_security_group" "main" {
  name   = "ecs-tasks-main-${var.module_id}"
  vpc_id = var.vpc.id

  ingress {
    from_port       = var.project.containers.main.ports.http.port
    to_port         = var.project.containers.main.ports.http.port
    protocol        = "tcp"
    security_groups = var.load_balancers.main.security_groups
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

# ECS Task Execution Role.
# https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/iam_role.
resource "aws_iam_role" "main" {
  name = var.iam_roles.main.name

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Principal = {
          Service = "ecs-tasks.amazonaws.com"
        }
      }
    ]
  })
}

# https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/iam_role_policy_attachment.
resource "aws_iam_role_policy_attachment" "main_policy_attachment" {
  role       = aws_iam_role.main.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AmazonECSTaskExecutionRolePolicy"
}

# Additional policy for ECR access if needed.
resource "aws_iam_role_policy" "main_policy" {
  name = "${var.iam_roles.main.name}-policy"
  role = aws_iam_role.main.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      # CloudWatch Logs.
      {
        Effect = "Allow"
        Action = [
          "logs:CreateLogGroup",
          "logs:CreateLogStream",
          "logs:PutLogEvents",
        ]
        Resource = [
          "arn:aws:logs:*:${var.aws_account_id}:log-group:/ecs/*:*",
          "arn:aws:logs:*:${var.aws_account_id}:log-group:/ecs/*:log-stream:*",
        ]
      },
      # ECR.
      {
        Effect = "Allow"
        Action = [
          "ecr:GetAuthorizationToken",
        ]
        Resource = "*"
      },
      # ECR.
      {
        Effect = "Allow"
        Action = [
          "ecr:BatchCheckLayerAvailability",
          "ecr:GetDownloadUrlForLayer",
          "ecr:BatchGetImage",
        ]
        Resource = "arn:aws:ecr:*:${var.aws_account_id}:repository/${var.deployment}/*"
      },
    ]
  })
}
