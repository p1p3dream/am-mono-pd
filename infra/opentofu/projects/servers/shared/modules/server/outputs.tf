output "certs" {
  sensitive = true
  value = {
    abodemine_main = {
      cert = aws_acm_certificate.abodemine_main
    }
  }
}

output "ecs_clusters" {
  value = {
    main_fargate = aws_ecs_cluster.main_fargate
  }
}

output "load_balancers" {
  value = {
    main = {
      load_balancer = aws_lb.main
      listeners = {
        http  = aws_lb_listener.main_http
        https = aws_lb_listener.main_https
      }
    }
  }
}
