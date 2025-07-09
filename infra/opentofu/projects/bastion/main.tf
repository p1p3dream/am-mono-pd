terraform {
  backend "s3" {
    region         = var.aws_region
    bucket         = var.s3_backend_bucket
    key            = var.s3_backend_keys.project
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

data "local_file" "config" {
  filename = "${path.root}/config.${var.deployment}.yaml"
}

locals {
  instance_id = "${var.project}-${random_string.instance_suffix.result}"
}

################################################################################
# MODULE: instance
################################################################################

resource "random_string" "instance_suffix" {
  length  = 8
  special = false
  upper   = false
}

module "instance" {
  source = "./modules/instance"

  aws_account_id = var.aws_account_id
  aws_region     = var.aws_region
  config         = yamldecode(data.local_file.config.content)
  deployment     = var.deployment
  module_id      = local.instance_id

  subnet_id = data.terraform_remote_state.main.outputs.vpc.public_subnets[1]

  tags                    = var.tags
  user_data_base64        = var.user_data_base64
  user_keypair_key_name   = local.instance_id
  user_keypair_public_key = var.user_keypair_public_key

  vpc = {
    id = data.terraform_remote_state.main.outputs.vpc.vpc_id
  }

  vpc_security_groups = {
    ssh = {
      ingress = {
        description = "SSH access"
        from_port   = 35857
        to_port     = 35857
        protocol    = "tcp"
        cidr_blocks = ["0.0.0.0/0"]
      }

      egress = {
        description = "All traffic"
        from_port   = 0
        to_port     = 0
        protocol    = "-1"
        cidr_blocks = ["0.0.0.0/0"]
      }
    }
  }

  vpc_security_group_ids = [
    data.terraform_remote_state.databases_os_alpha.outputs.security_groups.os_domain_users_sg.id,
    data.terraform_remote_state.databases_pg_alpha.outputs.security_groups.rds_pg_users_sg.id,
    data.terraform_remote_state.databases_pg_beta.outputs.security_groups.aurora_pg_users_sg.id,
    data.terraform_remote_state.databases_vk_alpha.outputs.security_groups.elasticache_vk_users_sg.id,
  ]
}
