# https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/iam_role.
# https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/iam_role_policy.

################################################################################
# ECS Task: main.
################################################################################

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
resource "aws_iam_role_policy_attachment" "main" {
  role       = aws_iam_role.main.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AmazonECSTaskExecutionRolePolicy"
}

# Additional policy for ECR access if needed.
resource "aws_iam_role_policy" "main" {
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
      {
        # Allow access to the locker DynamoDB table.
        Effect = "Allow"
        Action = [
          "dynamodb:GetItem",
          "dynamodb:PutItem",
          "dynamodb:UpdateItem",
          "dynamodb:Query",
          "dynamodb:Scan",
        ]
        Resource = var.dynamodb_tables.locker.table_arn
        Condition = {
          StringEquals = {
            "aws:ResourceTag/deployment" = var.deployment
          }
          "ForAllValues:StringLike" = {
            "dynamodb:LeadingKeys" = var.dynamodb_tables.locker.allowed_keys
          }
        }
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
      # S3.
      {
        # Allow management of partner data S3 buckets.
        Effect = "Allow"
        Action = [
          "s3:GetObject",
          "s3:ListBucket",
          "s3:PutObject",
        ]
        Resource = flatten([
          # Access to bucket root.
          [for arn in var.s3_buckets.partner_data : arn],
          # Access to bucket objects.
          [for arn in var.s3_buckets.partner_data : "${arn}/*"]
        ])
      },
      # SQS.
      {
        # Allow putting items in SQS queues.
        Effect = "Allow"
        Action = [
          "sqs:SendMessage",
        ]
        Resource = [
          for queue in var.sqs : queue.arn
        ]
      },
    ]
  })
}
