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

# CloudWatch

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
