output "partner_data_buckets" {
  description = "Map of all partner data created buckets."
  value = {
    for k, v in var.partner_data : k => aws_s3_bucket.partner_data[k]
  }
}

output "sqs_queues" {
  description = "Map of all SQS queues created."
  value = {
    main = aws_sqs_queue.main
  }
}
