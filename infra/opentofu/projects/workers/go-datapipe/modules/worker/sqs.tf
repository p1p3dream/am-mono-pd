# https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/sqs_queue.
resource "aws_sqs_queue" "main" {
  # .fifo suffix is required for FIFO queues.
  name = "${var.sqs_queue_name}.fifo"

  # https://docs.aws.amazon.com/AWSSimpleQueueService/latest/SQSDeveloperGuide/sqs-visibility-timeout.html.
  visibility_timeout_seconds = 180 # 3 minutes

  # This queue will receive messages from EventBridge every hour,
  # so we don't need to retain them for long.
  message_retention_seconds = 3600 # 1 hour

  max_message_size = 262144 # 256 KiB.

  # Process items immediately.
  delay_seconds = 0

  # We don't need to care about waiting for messages since
  # the queue length itself will trigger the ECS deployment.
  receive_wait_time_seconds = 0

  # https://docs.aws.amazon.com/AWSSimpleQueueService/latest/SQSDeveloperGuide/FIFO-queues-exactly-once-processing.html.
  # Use a fifo queue so we can have deduplication of messages
  # to ensure only one worker per job type is running at a time.
  fifo_queue                  = true
  content_based_deduplication = true
}
