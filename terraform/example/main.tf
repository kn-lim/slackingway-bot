locals {
  name                     = "slackingway-bot"
  account_id               = ""
  debug                    = "false"
  slack_signing_secret     = ""
  slack_oauth_token        = ""
  slack_history_channel_id = ""
  slack_output_channel_id  = ""
  admin_role_users         = ""

  # These non-empty .zip files are needed only when creating resources.
  # Run the build commands and zip the binary files.
  # The .zip file an be deleted/moved afterwards.
  endpoint_filename = "endpoint.zip"
  task_filename     = "task.zip"
}

module "slackingway-bot" {
  # https://github.com/kn-lim/slackingway-bot/tree/main/terraform
  source = "github.com/kn-lim/slackingway-bot//terraform"

  # Required

  account_id        = local.account_id
  endpoint_filename = local.endpoint_filename
  task_filename     = local.task_filename
  endpoint_environment_variables = {
    ADMIN_ROLE_USERS         = local.admin_role_users
    DEBUG                    = local.debug
    SLACK_HISTORY_CHANNEL_ID = local.slack_history_channel_id
    SLACK_OAUTH_TOKEN        = local.slack_oauth_token
    SLACK_OUTPUT_CHANNEL_ID  = local.slack_output_channel_id
    SLACK_SIGNING_SECRET     = local.slack_signing_secret
    TASK_FUNCTION_NAME       = "${local.name}-task"
  }
  task_environment_variables = {
    DEBUG                    = local.debug
    SLACK_HISTORY_CHANNEL_ID = local.slack_history_channel_id
    SLACK_OAUTH_TOKEN        = local.slack_oauth_token
    SLACK_OUTPUT_CHANNEL_ID  = local.slack_output_channel_id
    SLACK_SIGNING_SECRET     = local.slack_signing_secret
  }

  # Optional

  # name              = local.name
  # log_format        = "Text"
  # region            = "us-west-2"
  # retention_in_days = 3
  # runtime           = "provided.al2023"
  # endpoint_timeout  = 3
  # task_timeout      = 300
  # tags = {
  #   App = local.name
  # }
}

output "api_endpoint" {
  description = "The endpoint for the API Gateway"
  value       = module.slackingway-bot.api_endpoint
}
