################################################################################
# IAM
################################################################################

# IAM role in the DNS account for certificate validation
resource "aws_iam_role" "cert_validation" {
  provider = aws.management

  name = "cert-validation-${var.module_id}"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Principal = {
          AWS = "arn:aws:iam::${var.aws_account_id}:root"
        }
        Action = "sts:AssumeRole"
        Condition = {
          StringLike = {
            "aws:PrincipalArn" : [
              "arn:aws:iam::${var.aws_account_id}:role/aws-reserved/sso.amazonaws.com/*"
            ]
          }
        }
      }
    ]
  })
}

resource "aws_iam_role_policy" "cert_validation" {
  provider = aws.management
  name     = "cert-validation-policy-${var.module_id}"
  role     = aws_iam_role.cert_validation.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "route53:GetHostedZone",
          "route53:ChangeResourceRecordSets",
          "route53:ListResourceRecordSets"
        ]
        Resource = [for domain in var.domains : "arn:aws:route53:::hostedzone/${domain.zone_id}"]
      }
    ]
  })
}

################################################################################
# DOMAIN: abodemine_main
################################################################################

# Create the SSL certificate.
resource "aws_acm_certificate" "abodemine_main" {
  provider = aws.us_east_1

  domain_name               = var.domains.abodemine_main.name
  validation_method         = "DNS"
  subject_alternative_names = ["*.${var.domains.abodemine_main.name}"]

  lifecycle {
    create_before_destroy = true
  }
}

# Create the validation records.
# https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/route53_record.
resource "aws_route53_record" "abodemine_main" {
  provider = aws.management

  for_each = {
    for dvo in aws_acm_certificate.abodemine_main.domain_validation_options : dvo.domain_name => {
      name   = dvo.resource_record_name
      record = dvo.resource_record_value
      type   = dvo.resource_record_type
      domain = dvo.domain_name
    }
  }

  # MUST be enabled because multiple domains/SANs can require the same record.
  allow_overwrite = true

  zone_id = var.domains.abodemine_main.zone_id
  name    = each.value.name
  type    = each.value.type
  records = [each.value.record]
  ttl     = 60
}

# Validate the certificate.
resource "aws_acm_certificate_validation" "abodemine_main" {
  provider = aws.us_east_1

  certificate_arn         = aws_acm_certificate.abodemine_main.arn
  validation_record_fqdns = [for record in aws_route53_record.abodemine_main : record.fqdn]
}
