resource "aws_iam_openid_connect_provider" "github" {
  url = "https://token.actions.githubusercontent.com"

  client_id_list = ["sts.amazonaws.com"]

  thumbprint_list = [
    "6938fd4d98bab03faadb97b34396831e3780aea1"
  ]
}

resource "aws_iam_role" "github_actions" {
  name = var.role_name

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Principal = {
          Federated = aws_iam_openid_connect_provider.github.arn
        }
        Action = "sts:AssumeRoleWithWebIdentity"
        Condition = {
          StringLike = {
            "token.actions.githubusercontent.com:sub" = "repo:${var.github_org}/*:*"
          }
          StringEquals = {
            "token.actions.githubusercontent.com:aud" = "sts.amazonaws.com"
          }
        }
      },
      {
        # Allow AWS Reserved SSO AdministratorAccess roles to assume this role.
        # This is useful for local testing.
        Effect = "Allow"
        Principal = {
          AWS = "arn:aws:iam::${var.aws_account_id}:root"
        }
        Action = "sts:AssumeRole"
        Condition = {
          ArnLike = {
            "aws:PrincipalArn" = [
              "arn:aws:iam::${var.aws_account_id}:role/aws-reserved/sso.amazonaws.com/AWSReservedSSO_AdministratorAccess_*",
              "arn:aws:iam::${var.aws_account_id}:role/aws-reserved/sso.amazonaws.com/*/AWSReservedSSO_AdministratorAccess_*",
            ]
          }
        }
      }
    ]
  })
}

resource "aws_iam_role_policy" "github_actions" {
  name = "${var.role_name}-policy"
  role = aws_iam_role.github_actions.id

  policy = var.role_policy
}
