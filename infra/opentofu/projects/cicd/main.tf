terraform {
  backend "s3" {
    region         = var.aws_region
    bucket         = var.s3_backend_bucket
    key            = var.s3_backend_key
    dynamodb_table = var.s3_backend_table
  }

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.91.0"
    }

    random = {
      source  = "hashicorp/random"
      version = "~> 3.6.3"
    }
  }

  required_version = "~> 1.7"
}

locals {
  overlay_bucket_prefix = "${var.deployment}-mono-overlay"
}

################################################################################
# MODULE: build
################################################################################

resource "random_string" "build_suffix" {
  length  = 8
  special = false
  upper   = false
}

module "build" {
  source = "./modules/build"

  aws_region = var.aws_region



  s3_bucket_names = {
    build   = "mono-build-${random_string.build_suffix.result}",
    overlay = "${local.overlay_bucket_prefix}-${random_string.build_suffix.result}",
  }

  tags = var.tags
}

################################################################################
# MODULE: github_actions
################################################################################

resource "random_string" "github_actions_suffix" {
  length  = 8
  special = false
  upper   = false
}

module "github_actions" {
  source = "./modules/github_actions"

  aws_account_id = var.aws_account_id
  aws_region     = var.aws_region
  github_org     = var.github_organization_name
  role_name      = "github-actions-${random_string.github_actions_suffix.result}"

  role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        # Allow access to the DynamoDB table used for the Terraform backend.
        Effect = "Allow"
        Action = [
          "dynamodb:GetItem",
          "dynamodb:PutItem",
          "dynamodb:UpdateItem",
          "dynamodb:DeleteItem",
          "dynamodb:Query",
          "dynamodb:Scan",
        ]
        Resource = "arn:aws:dynamodb:*:${var.aws_account_id}:table/${var.s3_backend_table}"
        Condition = {
          StringEquals = {
            "aws:ResourceTag/deployment" = var.deployment
          }
          "ForAllValues:StringLike" = {
            "dynamodb:LeadingKeys" = ["${var.s3_backend_bucket}/${var.deployment}/*"]
          }
        }
      },
      {
        # Allow getting ECR authorization tokens.
        Effect = "Allow"
        Action = [
          "ecr:GetAuthorizationToken",
        ]
        Resource = "*"
      },
      {
        # Allow access to all ECR repositories in the deployment.
        Effect = "Allow"
        Action = [
          "ecr:PutImage",
          "ecr:InitiateLayerUpload",
          "ecr:UploadLayerPart",
          "ecr:CompleteLayerUpload",
          "ecr:BatchCheckLayerAvailability",
          "ecr:GetDownloadUrlForLayer",
          "ecr:BatchGetImage",
        ]
        Resource = "arn:aws:ecr:*:${var.aws_account_id}:repository/${var.deployment}/*"
      },
      {
        # 2025-02-20 (rbahiense):
        #   My goal was to restrict by account and deployment using resource
        #   matching and conditions, but I didn't find a way to implement it
        #   that satisfied the AWS API.
        #   Maybe the internal ECS API already limits based on region/account.
        #   Please only change it if requested, since it can be a time sink.
        Effect = "Allow"
        Action = [
          "ecs:DeregisterTaskDefinition",
          "ecs:DescribeTaskDefinition",
          "ecs:ListTaskDefinitions",
        ]
        Resource = [
          "*",
        ]
      },
      {
        # Allow managing existing ECS resources.
        Effect = "Allow"
        Action = [
          "ecs:DeleteService",
          "ecs:DescribeServices",
          "ecs:ListServices",
          "ecs:ListTagsForResource",
          "ecs:RunTask",
          "ecs:StartTask",
          "ecs:StopTask",
          "ecs:TagResource",
          "ecs:UpdateService",
        ]
        Resource = [
          "arn:aws:ecs:*:${var.aws_account_id}:service/*",
          "arn:aws:ecs:*:${var.aws_account_id}:task/*",
          "arn:aws:ecs:*:${var.aws_account_id}:task-definition/*",
        ]
        Condition = {
          StringEquals = {
            "ecs:ResourceTag/deployment" = var.deployment
          }
        }
      },
      {
        # Allow managing new ECS resources.
        Effect = "Allow"
        Action = [
          "ecs:CreateService",
          "ecs:RegisterTaskDefinition",
          "ecs:TagResource",
        ]
        Resource = [
          "arn:aws:ecs:*:${var.aws_account_id}:service/*",
          "arn:aws:ecs:*:${var.aws_account_id}:task/*",
          "arn:aws:ecs:*:${var.aws_account_id}:task-definition/*",
        ]
        Condition = {
          Null = {
            "ecs:ResourceTag/deployment" = "true"
          }
        }
      },
      {
        # Allow passing IAM roles to ECS tasks.
        Effect = "Allow"
        Action = [
          "iam:PassRole",
        ]
        Resource = [
          "arn:aws:iam::${var.aws_account_id}:role/ecs-task-*",
          "arn:aws:iam::${var.aws_account_id}:role/workers-go-datapipe-*",
        ]
        Condition = {
          StringEquals = {
            "iam:ResourceTag/deployment" = var.deployment
          }
        }
      },
      {
        # Allow access to the tofu S3 bucket.
        Effect = "Allow"
        Action = [
          "s3:GetObject",
          "s3:PutObject",
        ]
        Resource = "arn:aws:s3:::${var.s3_backend_bucket}/${var.deployment}/build/*"
      },
      {
        # Allow retrieval from deployment-specific overlay S3 buckets.
        Effect = "Allow"
        Action = [
          "s3:GetObject",
          "s3:ListBucket",
        ]
        Resource = "arn:aws:s3:::${local.overlay_bucket_prefix}-*"
      },
      {
        # Allow insertion into deployment-specific distribution S3 buckets.
        Effect = "Allow"
        Action = [
          "s3:ListBucket",
          "s3:PutObject",
        ]
        Resource = "arn:aws:s3:::${var.deployment}-www-*"
      },
      {
        # Allow access to the SSM parameters for the deployment.
        Effect = "Allow"
        Action = [
          "ssm:GetParameter",
          "ssm:GetParameters",
          "ssm:GetParametersByPath",
        ]
        Resource = "arn:aws:ssm:*:${var.aws_account_id}:parameter/${var.deployment}/*"
      },
    ]
  })

  tags = var.tags
}
