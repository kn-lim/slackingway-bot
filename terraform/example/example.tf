locals {
  name                     = "slackingway-bot"
  account_id               = ""
  slack_signing_secret     = ""
  slack_history_channel_id = ""
  slack_oauth_token        = ""
  endpoint_filename        = "endpoint.zip"
  task_filename            = "task.zip"
}

module "slackingway-bot" {
  source = "github.com/kn-lim/slackingway-bot//terraform"

  # Required

  account_id               = local.account_id
  endpoint_filename        = local.endpoint_filename
  slack_signing_secret     = local.slack_signing_secret
  slack_history_channel_id = local.slack_history_channel_id
  slack_oauth_token        = local.slack_oauth_token
  task_filename            = local.task_filename

  # Optional

  # endpoint_timeout = 3
  # function_name = local.name
  # log_format = "Text
  # region = "us-west-2"
  # retention_in_days = 3
  # runtime = "provided.al2023"
  # tags = {
  #   App = local.name
  # }
  # task_timeout = 3
}

output "api_endpoint" {
  description = "The endpoint for the API Gateway"
  value       = module.slackingway-bot.api_endpoint
}
