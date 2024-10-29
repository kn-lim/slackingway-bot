# slackingway-bot Terraform Module

To use this module, use the following as the source: `github.com/kn-lim/slackingway-bot//terraform`

Make sure to build the `endpoint` and `task` binaries in order for Terraform to create the resources. This will need to be done only when first creating the module.

<!-- BEGIN_TF_DOCS -->
## Example

```hcl
locals {
  name                     = "slackingway-bot"
  account_id               = ""
  slack_signing_secret     = ""
  slack_history_channel_id = ""
  slack_oauth_token        = ""

  # These non-empty .zip files are needed only when creating resources
  # Can be deleted/moved afterwards
  endpoint_filename = "endpoint.zip"
  task_filename     = "task.zip"
}

module "slackingway-bot" {
  # https://github.com/kn-lim/slackingway-bot/tree/main/terraform
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
```

## Providers

| Name | Version |
|------|---------|
| <a name="provider_aws"></a> [aws](#provider\_aws) | ~> 5.0 |

## Requirements

| Name | Version |
|------|---------|
| <a name="requirement_aws"></a> [aws](#requirement\_aws) | ~> 5.0 |

## Resources

| Name | Type |
|------|------|
| [aws_api_gateway_deployment.this](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/api_gateway_deployment) | resource |
| [aws_api_gateway_integration.this](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/api_gateway_integration) | resource |
| [aws_api_gateway_method.this](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/api_gateway_method) | resource |
| [aws_api_gateway_resource.this](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/api_gateway_resource) | resource |
| [aws_api_gateway_rest_api.this](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/api_gateway_rest_api) | resource |
| [aws_api_gateway_stage.this](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/api_gateway_stage) | resource |
| [aws_cloudwatch_log_group.endpoint](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/cloudwatch_log_group) | resource |
| [aws_cloudwatch_log_group.task](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/cloudwatch_log_group) | resource |
| [aws_iam_policy.lambda_logging](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/iam_policy) | resource |
| [aws_iam_role.endpoint](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/iam_role) | resource |
| [aws_iam_role.task](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/iam_role) | resource |
| [aws_iam_role_policy.invoke](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/iam_role_policy) | resource |
| [aws_iam_role_policy_attachment.lambda_logs_endpoint](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/iam_role_policy_attachment) | resource |
| [aws_iam_role_policy_attachment.lambda_logs_task](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/iam_role_policy_attachment) | resource |
| [aws_lambda_function.endpoint](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/lambda_function) | resource |
| [aws_lambda_function.task](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/lambda_function) | resource |
| [aws_lambda_permission.api_gateway](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/lambda_permission) | resource |
| [aws_iam_policy_document.assume_role](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/data-sources/iam_policy_document) | data source |
| [aws_iam_policy_document.lambda_logging](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/data-sources/iam_policy_document) | data source |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_account_id"></a> [account\_id](#input\_account\_id) | The AWS account ID | `string` | n/a | yes |
| <a name="input_endpoint_filename"></a> [endpoint\_filename](#input\_endpoint\_filename) | The filename of the endpoint Lambda function | `string` | n/a | yes |
| <a name="input_endpoint_timeout"></a> [endpoint\_timeout](#input\_endpoint\_timeout) | The timeout for the endpoint Lambda function | `number` | `3` | no |
| <a name="input_function_name"></a> [function\_name](#input\_function\_name) | The name of the Lambda function | `string` | `"slackingway-bot"` | no |
| <a name="input_log_format"></a> [log\_format](#input\_log\_format) | The log format for the CloudWatch logs | `string` | `"Text"` | no |
| <a name="input_region"></a> [region](#input\_region) | The region in which the resources will be created | `string` | `"us-west-2"` | no |
| <a name="input_retention_in_days"></a> [retention\_in\_days](#input\_retention\_in\_days) | The number of days to retain logs in CloudWatch | `number` | `3` | no |
| <a name="input_runtime"></a> [runtime](#input\_runtime) | The runtime for the Lambda functions | `string` | `"provided.al2023"` | no |
| <a name="input_slack_history_channel_id"></a> [slack\_history\_channel\_id](#input\_slack\_history\_channel\_id) | Slack channel ID for history | `string` | n/a | yes |
| <a name="input_slack_oauth_token"></a> [slack\_oauth\_token](#input\_slack\_oauth\_token) | Slack app OAuth token | `string` | n/a | yes |
| <a name="input_slack_signing_secret"></a> [slack\_signing\_secret](#input\_slack\_signing\_secret) | Slack app signing secret | `string` | n/a | yes |
| <a name="input_tags"></a> [tags](#input\_tags) | A map of tags to apply to the resources | `map(string)` | <pre>{<br/>  "App": "slackingway-bot"<br/>}</pre> | no |
| <a name="input_task_filename"></a> [task\_filename](#input\_task\_filename) | The filename of the task Lambda function | `string` | n/a | yes |
| <a name="input_task_timeout"></a> [task\_timeout](#input\_task\_timeout) | The timeout for the task Lambda function | `number` | `300` | no |

## Outputs

| Name | Description |
|------|-------------|
| <a name="output_api_endpoint"></a> [api\_endpoint](#output\_api\_endpoint) | The endpoint for the API Gateway |  
<!-- END_TF_DOCS -->
