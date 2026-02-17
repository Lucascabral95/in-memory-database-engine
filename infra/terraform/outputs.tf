output "ecr_repository_url" {
  description = "ECR repository URL used by CI/CD."
  value       = aws_ecr_repository.api.repository_url
}

output "ec2_instance_id" {
  description = "EC2 instance ID to target with SSM deploy commands."
  value       = aws_instance.api.id
}

output "ec2_public_ip" {
  description = "Public IP for direct API access."
  value       = aws_instance.api.public_ip
}

output "ssm_parameter_prefix" {
  description = "SSM prefix containing app environment variables."
  value       = local.ssm_parameter_prefix
}
