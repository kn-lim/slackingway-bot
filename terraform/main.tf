# Lambda

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

resource "aws_iam_role" "endpoint" {
  name               = "${var.function_name}-endpoint"
  assume_role_policy = data.aws_iam_policy_document.assume_role.json
  tags               = var.tags
}

resource "aws_lambda_function" "endpoint" {
  function_name = "${var.function_name}-endpoint"
  role          = aws_iam_role.endpoint.arn

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

  tags = var.tags

  depends_on = [aws_cloudwatch_log_group.endpoint]
}

resource "aws_iam_role" "task" {
  name               = "${var.function_name}-task"
  assume_role_policy = data.aws_iam_policy_document.assume_role.json
  tags               = var.tags
}

resource "aws_iam_role_policy" "invoke" {
  name = "InvokeEndpointLambdaFunction"
  role = aws_iam_role.task.name

  policy = jsonencode({
    Version = "2012-10-17",
    Statement = [
      {
        Effect   = "Allow",
        Action   = "lambda:InvokeFunction",
        Resource = aws_lambda_function.endpoint.arn,
      },
    ],
  })
}

resource "aws_lambda_function" "task" {
  function_name = "${var.function_name}-task"
  role          = aws_iam_role.task.arn

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
    log_format = "json"
  }

  tags = var.tags

  depends_on = [aws_cloudwatch_log_group.task]
}

# CloudWatch

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

data "aws_iam_policy_document" "lambda_logging" {
  statement {
    effect = "Allow"

    actions = [
      "logs:CreateLogGroup",
      "logs:CreateLogStream",
      "logs:PutLogEvents",
    ]

    resources = ["arn:aws:logs:*:*:*"]
  }
}

resource "aws_iam_policy" "lambda_logging" {
  name        = "lambda_logging"
  path        = "/"
  description = "IAM policy for logging from a lambda"
  policy      = data.aws_iam_policy_document.lambda_logging.json
}

resource "aws_iam_role_policy_attachment" "lambda_logs_endpoint" {
  role       = aws_iam_role.endpoint.name
  policy_arn = aws_iam_policy.lambda_logging.arn
}

resource "aws_iam_role_policy_attachment" "lambda_logs_task" {
  role       = aws_iam_role.task.name
  policy_arn = aws_iam_policy.lambda_logging.arn
}

# API Gateway

resource "aws_api_gateway_rest_api" "this" {
  name        = var.function_name
  description = "API Gateway for ${var.function_name}-endpoint"
  tags        = var.tags
}

resource "aws_api_gateway_resource" "this" {
  rest_api_id = aws_api_gateway_rest_api.this.id
  parent_id   = aws_api_gateway_rest_api.this.root_resource_id
  path_part   = "resource" # "{proxy+}"
}

resource "aws_api_gateway_method" "this" {
  rest_api_id   = aws_api_gateway_rest_api.this.id
  resource_id   = aws_api_gateway_resource.this.id
  http_method   = "GET"
  authorization = "NONE"
}

resource "aws_api_gateway_integration" "this" {
  rest_api_id = aws_api_gateway_rest_api.this.id
  resource_id = aws_api_gateway_resource.this.id
  http_method = aws_api_gateway_method.this.http_method

  integration_http_method = "POST"
  type                    = "AWS_PROXY"
  uri                     = aws_lambda_function.endpoint.invoke_arn
}

resource "aws_lambda_permission" "api_gateway" {
  statement_id  = "AllowExecutionFromAPIGateway"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.endpoint.function_name
  principal     = "apigateway.amazonaws.com"

  # More: http://docs.aws.amazon.com/apigateway/latest/developerguide/api-gateway-control-access-using-iam-policies-to-invoke-api.html
  source_arn = "arn:aws:execute-api:${var.region}:${var.account_id}:${aws_api_gateway_rest_api.this.id}/*/${aws_api_gateway_method.this.http_method}${aws_api_gateway_resource.this.path}"
}
