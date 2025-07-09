output "s3_buckets" {
  value = {
    partner_data = module.worker.partner_data_buckets
  }
}

output "sqs_queues" {
  value = module.worker.sqs_queues
}
