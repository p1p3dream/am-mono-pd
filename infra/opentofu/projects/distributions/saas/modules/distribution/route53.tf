# https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/route53_record.
resource "aws_route53_record" "abodemine_saas" {
  provider = aws.management

  zone_id = var.domains.abodemine_main.zone_id
  name    = var.domains.abodemine_saas.name
  type    = "A"

  alias {
    name                   = aws_cloudfront_distribution.main.domain_name
    zone_id                = aws_cloudfront_distribution.main.hosted_zone_id
    evaluate_target_health = true
  }
}
