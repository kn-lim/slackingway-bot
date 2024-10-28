resource "aws_cloudwatch_log_group" "endpoint" {
  name              = "/aws/lambda/${aws_lambda_function.endpoint.function_name}"
  retention_in_days = var.retention_in_days
  tags              = var.tags
}

resource "aws_cloudwatch_log_group" "task" {
  name              = "/aws/lambda/${aws_lambda_function.task.function_name}"
  retention_in_days = var.retention_in_days
  tags              = var.tags
}
