<p align="center">
  <img width="100" style="border-radius: 50%" src="https://raw.githubusercontent.com/kn-lim/slackingway-bot/main/images/slackingway.png"></img>
  <br>
  <i>I'm a</i> 🤖<i>!</i>
</p>

# slackingway-bot

![Go](https://img.shields.io/github/go-mod/go-version/kn-lim/slackingway-bot)
![GitHub Workflow Status - Build](https://img.shields.io/github/actions/workflow/status/kn-lim/slackingway-bot/build.yaml)
![GitHub Workflow Status - Test](https://img.shields.io/github/actions/workflow/status/kn-lim/slackingway-bot/test.yaml?label=tests)
[![Coverage Status](https://coveralls.io/repos/github/kn-lim/slackingway-bot/badge.svg?branch=main)](https://coveralls.io/github/kn-lim/slackingway-bot?branch=main)
[![Go Report Card](https://goreportcard.com/badge/github.com/kn-lim/slackingway-bot)](https://goreportcard.com/report/github.com/kn-lim/slackingway-bot)
![License](https://img.shields.io/github/license/kn-lim/slackingway-bot)

A personal Slack bot to handle miscellaneous tasks hosted on AWS Lambda.

## Packages Used

- [aws-lambda-go](https://github.com/aws/aws-lambda-go/)
- [aws-sdk-go-v2](https://github.com/aws/aws-sdk-go-v2)
- [chattingway](https://github.com/kn-lim/chattingway)
- [mock](https://github.com/uber-go/mock)
- [slack-go](https://github.com/slack-go/slack)
- [testify](https://github.com/stretchr/testify)

# Using the Slack Bot

## Slack Slash Commands

| Command | Description |
| - | - |
| `/coinflip` | Flips a coin |
| `/delayed-ping` | Ping with a delay |
| `/echo` | Opens a Slack modal to echo a text to the output channel |
| `/menu` | Opens a Slack modal to select options and sends the result to the output channel |
| `/ping` | Ping |
| `/roll` | Rolls a dice with modifiers |

## How to Build

From the project home directory: 

- **Endpoint Function**: `CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -tags lambda.norpc -o binary/bootstrap ./cmd/endpoint/`
- **Task Function**: `CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -tags lambda.norpc -o binary/bootstrap ./cmd/task/`

Zip the `bootstrap` binaries and upload it to the Lambda functions.

## Environment Variables

### Endpoint Lambda Function

| Name | Description |
| - | - |
| `DEBUG` | Enable debug mode |
| `TASK_FUNCTION_NAME` | Name of the Task Lambda Function |
| `SLACK_SIGNING_SECRET` | Slack App's Signing Secret |
| `SLACK_OAUTH_TOKEN` | Slack App's OAuth Token |
| `SLACK_HISTORY_CHANNEL_ID` | Slackingway's History Channel ID |
| `SLACK_OUTPUT_CHANNEL_ID` | Slackingway's Output Channel ID |
| `ADMIN_ROLE_USERS` | Comma-delimited string of Slack User IDs with admin roles |

### Task Lambda Function

| Name | Description |
| - | - |
| `DEBUG` | Enable debug mode |
| `SLACK_OAUTH_TOKEN` | Slack App's OAuth Token |
| `SLACK_HISTORY_CHANNEL_ID` | Slackingway's History Channel ID |
| `SLACK_OUTPUT_CHANNEL_ID` | Slackingway's Output Channel ID |

## AWS Setup

To quickly spin up **slackingway-bot** on AWS, use the [Terraform module](https://github.com/kn-lim/chattingway-terraform/).

1. Create the **endpoint** Lambda function on AWS. 
    - For the `Runtime`, select `Amazon Linux 2023`.
    - For the `Architecture`, select `x86_64`.
2. Add an API Gateway triger to the **endpoint** Lambda function.
    - Use the following settings:
      - **Intent**: Create a new API
      - **API type**: REST API
      - **Security**: Open
3. Create the **task** Lambda function on AWS. 
    - For the `Runtime`, select `Amazon Linux 2023`.
    - For the `Architecture`, select `x86_64`.
4. Build the **endpoint** and **task** binaries.
5. Archive the `bootstrap` binaries in .zip files and upload it to the Lambda functions.
6. In the `Configuration` tab, add in the required environment variables to the Lambda functions.
7. Change the `Timeout` of the **task** Lambda function to a value greater than 3 seconds.
    - The `Timeout` of the **endpoint** Lambda function can stay as 3 seconds to follow Slack's requirements.

## Slack Setup

### Slash Commands

Get the **endpoint** Lambda API Gateway triggers's `API endpoint` and add it to the Slack apps's `Request URL` in each Slack Slash Command in the Slack API page.

### OAuth & Permissions

#### OAuth Tokens

Save the `Bot User OAuth Token` as the `SLACK_OAUTH_TOKEN` environment variable in the **task** Lambda function.

#### Scopes

Enable the following `Bot Token Scopes`:

- `channels:history`
- `chat:write`
- `chat:write.customize`
- `commands`
- `im:history`
- `users.profile:read`
- `users:read`

### Event Subscriptions

#### Enable Events

Get the **endpoint** Lambda API Gateway triggers's `API endpoint` and add it to the Slack apps's `Request URL`. It should be verified after a second.

#### Subscribe to Bot Events

Add the following bot user events:

- `app_home_opened`
