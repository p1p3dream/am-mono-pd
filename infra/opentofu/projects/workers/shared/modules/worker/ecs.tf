# https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/ecs_cluster.html.
resource "aws_ecs_cluster" "main_fargate" {
  name = var.ecs_clusters.main_fargate.name

  setting {
    name  = var.ecs_clusters.main_fargate.setting.name
    value = var.ecs_clusters.main_fargate.setting.value
  }
}

# https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/ecs_cluster_capacity_providers.
resource "aws_ecs_cluster_capacity_providers" "main_fargate" {
  cluster_name = aws_ecs_cluster.main_fargate.name

  capacity_providers = ["FARGATE", "FARGATE_SPOT"]

  default_capacity_provider_strategy {
    base              = 1
    weight            = 100
    capacity_provider = "FARGATE"
  }
}
