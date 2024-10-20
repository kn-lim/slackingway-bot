variable "region" {
  description = "The region in which the resources will be created"
  type        = string
  default     = "us-west-2"
}

variable "slack_signing_secret" {
  description = "Slack app signing secret"
  type        = string
  sensitive   = true
}

variable "function_name" {
  description = "The name of the Lambda function"
  type        = string
  default     = "slackingway-bot"
}

variable "runtime" {
  description = "The runtime for the Lambda functions"
  type        = string
  default     = "provided.al2023"
}

variable "endpoint_timeout" {
  description = "The timeout for the endpoint Lambda function"
  type        = number
  default     = 3
}

variable "task_timeout" {
  description = "The timeout for the task Lambda function"
  type        = number
  default     = 300
}

variable "retention_in_days" {
  description = "The number of days to retain logs in CloudWatch"
  type        = number
  default     = 3
}
