module "slackingway-bot" {
  source = "github.com/kn-lim/slackingway-bot//terraform"

  # Required

  account_id               = ""
  slack_signing_secret     = ""
  slack_history_channel_id = ""
  slack_oauth_token        = ""

  endpoint_filename = "endpoint.zip"
  task_filename     = "task.zip"

  # Optional

  # endpoint_timeout = 3
  # function_name = "slackingway-bot"
  # log_format = "Text
  # region = "us-west-2"
  # retention_in_days = 3
  # runtime = "provided.al2023"
  # tags = {
  #   App = "slackingway-bot"
  # }
  # task_timeout = 3
}
