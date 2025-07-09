output "ecs_clusters" {
  value = {
    main_fargate = aws_ecs_cluster.main_fargate
  }
}
