variable "region" {
  description = "The region in which the resources will be created"
  type        = string
  default     = "us-west-2"
}

variable "account_id" {
  description = "The AWS account ID"
  type        = string
}

variable "slack_signing_secret" {
  description = "Slack app signing secret"
  type        = string
  sensitive   = true
}

variable "slack_history_channel_id" {
  description = "Slack channel ID for history"
  type        = string
}

variable "slack_oauth_token" {
  description = "Slack app OAuth token"
  type        = string
  sensitive   = true
}

variable "endpoint_filename" {
  description = "The filename of the endpoint Lambda function"
  type        = string
}

variable "task_filename" {
  description = "The filename of the task Lambda function"
  type        = string
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

variable "tags" {
  description = "A map of tags to apply to the resources"
  type        = map(string)
  default = {
    App = "slackingway-bot"
  }
}
