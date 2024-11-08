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

variable "slack_output_channel_id" {
  description = "Slack channel ID for output"
  type        = string
}

variable "filename" {
  description = "The filename to upload to the Lambda function"
  type        = string
}

variable "name" {
  description = "The name of the resources"
  type        = string
  default     = "slackingway-bot"
}

variable "runtime" {
  description = "The runtime for the Lambda functions"
  type        = string
  default     = "provided.al2023"
}

variable "timeout" {
  description = "The timeout for the Lambda function"
  type        = number
  default     = 300
}

variable "log_format" {
  description = "The log format for the CloudWatch logs"
  type        = string
  default     = "Text"
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

variable "debug" {
  description = "Enable debug mode"
  type        = string
  default     = false
}