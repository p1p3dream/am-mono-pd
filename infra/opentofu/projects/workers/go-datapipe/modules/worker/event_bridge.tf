# https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/scheduler_schedule.

resource "aws_scheduler_schedule" "sqs_scheduler" {
  for_each = var.partner_schedule

  name       = each.value.schedule_name
  group_name = "default"

  flexible_time_window {
    mode = "OFF"
  }

  schedule_expression = each.value.schedule_expression

  target {
    arn      = aws_sqs_queue.main.arn
    role_arn = aws_iam_role.eventbridge_sqs_scheduler.arn

    input = each.value.sqs_message_body

    retry_policy {
      maximum_retry_attempts = 3
    }

    sqs_parameters {
      message_group_id = each.value.sqs_message_group_id
    }
  }
}
