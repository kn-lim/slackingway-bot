resource "aws_lambda_function" "endpoint" {
  filename = var.endpoint_filename
  function_name = "${var.function_name}-endpoint"
  role          = aws_iam_role.endpoint.arn
  handler = "hello.handler" # Not used
  runtime = var.runtime
  timeout = var.endpoint_timeout

  environment {
    variables = {
      TASK_FUNCTION_NAME   = "${var.function_name}-task"
      SLACK_SIGNING_SECRET = var.slack_signing_secret
    }
  }

  logging_config {
    log_group  = aws_cloudwatch_log_group.endpoint.name
    log_format = var.log_format
  }

  tags = var.tags

  depends_on = [aws_cloudwatch_log_group.endpoint]
}

resource "aws_lambda_function" "task" {
  filename = var.task_filename
  function_name = "${var.function_name}-task"
  role          = aws_iam_role.task.arn
  handler = "hello.handler" # Not used
  runtime = var.runtime
  timeout = var.task_timeout

  environment {
    variables = {
      SLACK_HISTORY_CHANNEL_ID = var.slack_history_channel_id
      SLACK_OAUTH_TOKEN        = var.slack_oauth_token
    }
  }

  logging_config {
    log_group  = aws_cloudwatch_log_group.task.name
    log_format = var.log_format
  }

  tags = var.tags

  depends_on = [aws_cloudwatch_log_group.task]
}

resource "aws_lambda_permission" "api_gateway" {
  statement_id  = "AllowExecutionFromAPIGateway"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.endpoint.function_name
  principal     = "apigateway.amazonaws.com"

  # More: http://docs.aws.amazon.com/apigateway/latest/developerguide/api-gateway-control-access-using-iam-policies-to-invoke-api.html
  source_arn = "arn:aws:execute-api:${var.region}:${var.account_id}:${aws_api_gateway_rest_api.this.id}/*/${aws_api_gateway_method.this.http_method}${aws_api_gateway_resource.this.path}"
}
