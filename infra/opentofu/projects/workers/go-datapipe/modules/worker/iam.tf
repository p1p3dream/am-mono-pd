# https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/iam_role.
# https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/iam_role_policy.

################################################################################
# EventBridge SQS Scheduler.
################################################################################

resource "aws_iam_role" "eventbridge_sqs_scheduler" {
  name = var.iam_roles.eventbridge_sqs_scheduler.name

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Principal = {
          Service = "scheduler.amazonaws.com"
        }
      }
    ]
  })
}

resource "aws_iam_role_policy" "eventbridge_sqs_scheduler_policy" {
  name = "${var.iam_roles.eventbridge_sqs_scheduler.name}-policy"
  role = aws_iam_role.eventbridge_sqs_scheduler.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "sqs:SendMessage"
        ]
        Resource = [aws_sqs_queue.main.arn]
      }
    ]
  })
}
