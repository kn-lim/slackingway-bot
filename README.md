<p align="center">
  <img width="100" style="border-radius: 50%" src="https://raw.githubusercontent.com/kn-lim/slackingway-bot/main/images/slackingway.png"></img>
  <br>
  <i>I'm a</i> 🤖<i>!</i>
</p>

# slackingway-bot

![Go](https://img.shields.io/github/go-mod/go-version/kn-lim/slackingway-bot)
![GitHub Workflow Status - Build](https://img.shields.io/github/actions/workflow/status/kn-lim/slackingway-bot/build.yaml)
[![Go Report Card](https://goreportcard.com/badge/github.com/kn-lim/slackingway-bot)](https://goreportcard.com/report/github.com/kn-lim/slackingway-bot)
![License](https://img.shields.io/github/license/kn-lim/slackingway-bot)

A personal Slack bot to handle miscellaneous tasks hosted on AWS Lambda.

## Packages Used

- [aws-lambda-go](https://github.com/aws/aws-lambda-go/)
- [slack-go](https://github.com/slack-go/slack)

# Using the Slack Bot

## How to Build

From the project home directory: 

- **Endpoint Function**: `CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -tags lambda.norpc -o binary/bootstrap ./cmd/endpoint/`
- **Task Function**: `CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -tags lambda.norpc -o binary/bootstrap ./cmd/task/`

### Build Environment Variables

| Name | Description |
| - | - |
| `AWS_REGION` | AWS Region of the Task Lambda Function |

## Environment Variables

### Endpoint Lambda Function

#### AWS

| Name | Description |
| - | - |
| `TASK_FUNCTION_NAME` | Name of the Task Lambda Function |

#### Slack

| Name | Description |
| - | - |
| `SLACK_SIGNING_SECRET` | Slack App Signing Secret |

### Task Lambda Function

| Name | Description |
| - | - |
| `SLACK_SIGNING_SECRET` | Slack App Signing Secret |

## AWS Setup

1. Create an **endpoint** Lambda function on AWS. 
    - For the `Runtime`, select `Provide your own bootstrap on Amazon Linux 2` under `Custom runtime`.
    - For the `Architecture`, select `x86_64`.
    - Under `Advanced Settings`, select:
        - `Enable function URL`
          - Auth type: `NONE`
          - Invoke mode: `BUFFERED (default)`
          - Enable `Configure cross-origin resource sharing (CORS)`
2. Add an API Gateway triger to the **endpoint** Lambda function.
    - Use the following settings:
      - **Intent**: Create a new API
      - **API type**: REST API
      - **Security**: Open
3. Create a **task** Lambda function on AWS. 
    - For the `Runtime`, select `Provide your own bootstrap on Amazon Linux 2` under `Custom runtime`.
    - For the `Architecture`, select `x86_64`.
4. Archive the `bootstrap` binary in a .zip file and upload it to the Lambda functions.
5. In the `Configuration` tab, add in the required environment variables to the Lambda functions.
6. Give the role the **endpoint** Lambda function is using permission to run the **task** Lambda function.
7. Give the role the **task** Lambda function is using permission to access the AWS resources it will need.
8. Change the `Timeout` of the **task** Lambda function to a value greater than 3 seconds.
    - The `Timeout` of the **endpoint** Lambda function can stay as 3 seconds to follow Slack's requirements.
9. Get the **endpoint** Lambda API Gateway triggers's `API endpoint` and add it to the Slack apps's `Request URL` in each Slack Slash Command in the [Slack Apps](https://api.slack.com/apps/) page.
