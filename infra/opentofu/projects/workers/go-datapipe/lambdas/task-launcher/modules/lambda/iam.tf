# https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/iam_role.
# https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/iam_role_policy.

################################################################################
# Lambda
################################################################################

resource "aws_iam_role" "lambda" {
  name = var.iam_roles.lambda.name

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Principal = {
          Service = "lambda.amazonaws.com"
        }
      }
    ]
  })
}

resource "aws_iam_role_policy" "lambda_policy" {
  name = "${var.iam_roles.lambda.name}-policy"
  role = aws_iam_role.lambda.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      # Allow creating CloudWatch log streams and inserting events.
      {
        Effect = "Allow"
        Action = [
          "logs:CreateLogStream",
          "logs:PutLogEvents"
        ]
        Resource = [
          "arn:aws:logs:${var.aws_region}:${var.aws_account_id}:*"
        ]
      },
      {
        # Allow reading the locker DynamoDB table.
        Effect = "Allow"
        Action = [
          "dynamodb:GetItem",
          "dynamodb:Query",
          "dynamodb:Scan",
        ]
        Resource = var.dynamodb.locker.table_arn
        Condition = {
          StringEquals = {
            "aws:ResourceTag/deployment" = var.deployment
          }
          "ForAllValues:StringLike" = {
            "dynamodb:LeadingKeys" = var.dynamodb.locker.allowed_keys
          }
        }
      },
      # Allow listing ECS task definitions.
      {
        Effect = "Allow"
        Action = [
          "ecs:DescribeTaskDefinition",
          "ecs:ListTaskDefinitions",
        ]
        Resource = [
          "*",
        ]
      },
      # Allow running ECS tasks on the given cluster.
      {
        Effect = "Allow"
        Action = [
          "ecs:RunTask",
        ]
        Resource = var.lambda.allow_run_task_on
      },
      {
        # Allow passing IAM roles to ECS tasks.
        Effect = "Allow"
        Action = [
          "iam:PassRole",
        ]
        Resource = var.lambda.allow_pass_role_on
      },
      # Allow fetching SQS queue attributes and managing messages.
      {
        Effect = "Allow"
        Action = [
          "sqs:ChangeMessageVisibility",
          "sqs:GetQueueAttributes",
          "sqs:DeleteMessage",
          "sqs:ReceiveMessage",
        ]
        Resource = [
          var.sqs.main.arn
        ]
      },
    ]
  })
}
