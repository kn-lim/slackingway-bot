output "api_endpoint" {
  description = "The endpoint for the API Gateway"
  value       = "${aws_api_gateway_stage.this.invoke_url}/endpoint"
}
