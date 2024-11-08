locals {
  name                     = "slackingway-bot"
  account_id               = ""
  slack_signing_secret     = ""
  slack_oauth_token        = ""
  slack_history_channel_id = ""
  slack_output_channel_id  = ""

  # This non-empty .zip file is needed only when creating resources.
  # Run the build command and zip the binary file.
  # The .zip file an be deleted/moved afterwards.
  filename = "bootstrap.zip"
}

module "slackingway-bot" {
  # https://github.com/kn-lim/slackingway-bot/tree/main/terraform
  source = "github.com/kn-lim/slackingway-bot//terraform"

  # Required

  account_id               = local.account_id
  filename                 = local.filename
  slack_oauth_token        = local.slack_oauth_token
  slack_signing_secret     = local.slack_signing_secret
  slack_history_channel_id = local.slack_history_channel_id
  slack_output_channel_id  = local.slack_output_channel_id

  # Optional

  # debug             = "false" 
  # name              = local.name
  # log_format        = "Text"
  # region            = "us-west-2"
  # retention_in_days = 3
  # runtime           = "provided.al2023"
  # tags = {
  #   App = local.name
  # }
  # timeout = 300
}

output "api_endpoint" {
  description = "The endpoint for the API Gateway"
  value       = module.slackingway-bot.api_endpoint
}
