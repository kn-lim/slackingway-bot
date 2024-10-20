data "aws_iam_policy_document" "assume_role" {
  statement {
    effect = "Allow"

    principals {
      type        = "Service"
      identifiers = ["lambda.amazonaws.com"]
    }

    actions = ["sts:AssumeRole"]
  }
}

resource "aws_iam_role" "lambda" {
  name               = "lambda"
  assume_role_policy = data.aws_iam_policy_document.assume_role.json
}

resource "aws_lambda_function" "endpoint" {
  function_name = "${var.function_name}-endpoint"
  role          = aws_iam_role.lambda.arn

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
    log_format = "json"
  }

  depends_on = [aws_cloudwatch_log_group.endpoint]
}

resource "aws_lambda_function" "task" {
  function_name = "${var.function_name}-task"
  role          = aws_iam_role.lambda.arn

  runtime = var.runtime
  timeout = var.task_timeout

  environment {
    variables = {
      SLACK_SIGNING_SECRET = var.slack_signing_secret
    }
  }

  logging_config {
    log_group  = aws_cloudwatch_log_group.task.name
    log_format = "json"
  }

  depends_on = [aws_cloudwatch_log_group.task]
}

resource "aws_cloudwatch_log_group" "endpoint" {
  name              = "/aws/lambda/${aws_lambda_function.endpoint.function_name}"
  retention_in_days = var.retention_in_days
}

resource "aws_cloudwatch_log_group" "task" {
  name              = "/aws/lambda/${aws_lambda_function.task.function_name}"
  retention_in_days = var.retention_in_days
}
