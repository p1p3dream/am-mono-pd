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
        # Allow reading the secure_download DynamoDB table.
        Effect = "Allow"
        Action = [
          "dynamodb:GetItem",
          "dynamodb:Query",
          "dynamodb:Scan",
        ]
        Resource = var.dynamodb_tables.secure_download.arn
      },
      {
        Effect = "Allow"
        Action = [
          "s3:GetObject",
          "s3:ListBucket",
        ]
        Resource = [
          var.s3_buckets.secure_download.arn,
          "${var.s3_buckets.secure_download.arn}/*"
        ]
      },
    ]
  })
}
