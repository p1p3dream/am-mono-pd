# https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/opensearch_domain

# Can only be created once per account.
# Check with aws iam get-role --role-name AWSServiceRoleForAmazonOpenSearchService.
# resource "aws_iam_service_linked_role" "opensearch" {
#   aws_service_name = "opensearchservice.amazonaws.com"
# }

data "aws_iam_policy_document" "os" {
  statement {
    effect = "Allow"

    principals {
      type        = "AWS"
      identifiers = ["*"]
    }

    actions   = ["es:*"]
    resources = ["arn:aws:es:${var.aws_region}:${var.aws_account_id}:domain/${var.os_domain.domain_name}/*"]
  }
}

resource "aws_opensearch_domain" "os" {
  domain_name    = var.os_domain.domain_name
  engine_version = var.os_domain.engine_version

  cluster_config {
    instance_type  = var.os_domain.cluster_config.instance_type
    instance_count = var.os_domain.cluster_config.instance_count

    dedicated_master_enabled = var.os_domain.cluster_config.dedicated_master_enabled
    dedicated_master_type    = var.os_domain.cluster_config.dedicated_master_type
    dedicated_master_count   = var.os_domain.cluster_config.dedicated_master_count

    warm_enabled           = var.os_domain.cluster_config.warm_enabled
    zone_awareness_enabled = var.os_domain.cluster_config.zone_awareness_enabled

    dynamic "zone_awareness_config" {
      for_each = var.os_domain.cluster_config.zone_awareness_enabled ? [1] : []
      content {
        availability_zone_count = var.os_domain.cluster_config.zone_awareness_config.availability_zone_count
      }
    }
  }

  ebs_options {
    ebs_enabled = var.os_domain.ebs_options.ebs_enabled
    volume_size = var.os_domain.ebs_options.volume_size
    volume_type = var.os_domain.ebs_options.volume_type
    iops        = var.os_domain.ebs_options.iops
    throughput  = var.os_domain.ebs_options.throughput
  }

  encrypt_at_rest {
    enabled = var.os_domain.encrypt_at_rest.enabled
  }

  node_to_node_encryption {
    enabled = var.os_domain.node_to_node_encryption.enabled
  }

  domain_endpoint_options {
    enforce_https       = var.os_domain.domain_endpoint_options.enforce_https
    tls_security_policy = var.os_domain.domain_endpoint_options.tls_security_policy
  }

  advanced_security_options {
    enabled                        = var.os_domain.advanced_security_options.enabled
    internal_user_database_enabled = var.os_domain.advanced_security_options.internal_user_database_enabled

    master_user_options {
      master_user_name     = var.os_domain.advanced_security_options.master_user_options.master_user_name
      master_user_password = var.os_domain.advanced_security_options.master_user_options.master_user_password
    }
  }

  vpc_options {
    subnet_ids = var.vpc.database_subnets
    security_group_ids = [
      aws_security_group.os_domain_sg.id,
    ]
  }

  access_policies = data.aws_iam_policy_document.os.json
}
