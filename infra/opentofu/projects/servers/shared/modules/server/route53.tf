# https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/route53_record.
resource "aws_route53_record" "abodemine_api" {
  provider = aws.management

  zone_id = var.domains.abodemine_main.zone_id
  name    = var.domains.abodemine_api.name
  type    = "A"

  alias {
    name                   = aws_lb.main.dns_name
    zone_id                = aws_lb.main.zone_id
    evaluate_target_health = true
  }
}
