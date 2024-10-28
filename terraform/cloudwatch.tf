resource "aws_cloudwatch_log_group" "endpoint" {
  name              = "/aws/lambda/${var.function_name}-endpoint"
  retention_in_days = var.retention_in_days
  tags              = var.tags
}

resource "aws_cloudwatch_log_group" "task" {
  name              = "/aws/lambda/${var.function_name}-task"
  retention_in_days = var.retention_in_days
  tags              = var.tags
}
