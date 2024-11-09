output "api_endpoint" {
  description = "The endpoint for the API Gateway"
  value       = "${aws_api_gateway_stage.this.invoke_url}/${aws_api_gateway_resource.this.path_part}"
}

output "endpoint_function_arn" {
  description = "The ARN of the Endpoint Lambda function"
  value       = "${aws_lambda_function.endpoint.arn}"
}

output "task_function_arn" {
  description = "The ARN of the Task Lambda function"
  value       = "${aws_lambda_function.task.arn}"
}