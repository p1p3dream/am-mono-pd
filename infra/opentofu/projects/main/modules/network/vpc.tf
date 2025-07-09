# https://registry.terraform.io/modules/terraform-aws-modules/vpc/aws/latest.
# https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/opensearch_vpc_endpoint.

# Filter out local zones, which we don't need/want.
data "aws_availability_zones" "available" {
  filter {
    name   = "opt-in-status"
    values = ["opt-in-not-required"]
  }
}

module "vpc" {
  source  = "terraform-aws-modules/vpc/aws"
  version = var.vpc.version

  name       = var.vpc_name
  create_vpc = true

  cidr = "10.0.0.0/16"
  azs  = slice(data.aws_availability_zones.available.names, 0, 3)

  database_subnets    = ["10.0.201.0/24", "10.0.202.0/24", "10.0.203.0/24"]
  elasticache_subnets = ["10.0.211.0/24", "10.0.212.0/24", "10.0.213.0/24"]
  private_subnets     = ["10.0.101.0/24", "10.0.102.0/24", "10.0.103.0/24"]
  public_subnets      = ["10.0.1.0/24", "10.0.2.0/24", "10.0.3.0/24"]

  # default_security_group_egress = []
  # default_security_group_ingress = []

  create_database_internet_gateway_route = var.vpc.create_database_internet_gateway_route
  create_database_subnet_group           = true
  create_database_subnet_route_table     = var.vpc.create_database_subnet_route_table
  create_elasticache_subnet_group        = true
  enable_dns_hostnames                   = true
  enable_dns_support                     = true
  enable_nat_gateway                     = true
  map_public_ip_on_launch                = false # Always false because workers MUST NOT have public IPs.
  single_nat_gateway                     = true  # Set to false for HA.
}

# Based on https://github.com/terraform-aws-modules/terraform-aws-vpc/blob/master/examples/complete/main.tf.
module "vpc_endpoints" {
  source  = "terraform-aws-modules/vpc/aws//modules/vpc-endpoints"
  version = var.vpc.version

  vpc_id = module.vpc.vpc_id

  create_security_group      = true
  security_group_name_prefix = "${var.vpc_name}-vpc-endpoints-"
  security_group_description = "VPC endpoint security group."
  security_group_rules = {
    ingress_https = {
      description = "HTTPS from VPC."
      cidr_blocks = [module.vpc.vpc_cidr_block]
    }
  }

  endpoints = {
    dynamodb = {
      service         = "dynamodb"
      service_type    = "Gateway"
      route_table_ids = flatten([module.vpc.intra_route_table_ids, module.vpc.private_route_table_ids, module.vpc.public_route_table_ids])
    },

    # ecr_api = {
    #   service             = "ecr.api"
    #   private_dns_enabled = true
    #   subnet_ids          = module.vpc.private_subnets
    # },

    # ecr_dkr = {
    #   service             = "ecr.dkr"
    #   private_dns_enabled = true
    #   subnet_ids          = module.vpc.private_subnets
    # },

    # ecs = {
    #   service             = "ecs"
    #   private_dns_enabled = true
    #   subnet_ids          = module.vpc.private_subnets
    # },

    # ecs_telemetry = {
    #   create              = false
    #   service             = "ecs-telemetry"
    #   private_dns_enabled = true
    #   subnet_ids          = module.vpc.private_subnets
    # },

    # logs = {
    #   service             = "logs"
    #   private_dns_enabled = true
    #   subnet_ids          = module.vpc.private_subnets
    # },

    # Kept for future need.
    # rds = {
    #   service             = "rds"
    #   private_dns_enabled = true
    #   subnet_ids          = module.vpc.private_subnets
    # },

    # redshift = {
    #   service             = "redshift"
    #   private_dns_enabled = true
    #   subnet_ids          = module.vpc.private_subnets
    # }

    # redshift_data = {
    #   service             = "redshift-data"
    #   private_dns_enabled = true
    #   subnet_ids          = module.vpc.private_subnets
    # }

    s3 = {
      service             = "s3"
      service_type        = "Gateway"
      private_dns_enabled = true
      route_table_ids     = flatten([module.vpc.intra_route_table_ids, module.vpc.private_route_table_ids, module.vpc.public_route_table_ids])
    },
  }
}

# # OpenSearch VPC endpoint using the same config as the VPC endpoints module.
# # This is necessary because the VPC endpoints module does not support OpenSearch yet (2025-01-18).
# resource "aws_opensearch_vpc_endpoint" "main" {
#   domain_arn         = aws_opensearch_domain.domain.arn
#   vpc_id             = module.vpc.vpc_id
#   subnet_ids         = module.vpc.private_subnets
#   security_group_ids = [module.vpc_endpoints.security_group_id]
# }
