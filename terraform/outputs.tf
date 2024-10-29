output "api_endpoint" {
  description = "The endpoint for the API Gateway"
  value       = "https://${aws_api_gateway_rest_api.this.id}.execute-api.${var.region}.amazonaws.com/default/${aws_api_gateway_rest_api.this.name}"
}
