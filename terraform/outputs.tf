output "ec2_public_ip" {
  description = "Public IP address of the EC2 instance"
  value       = aws_instance.app.public_ip
}

output "ec2_public_dns" {
  description = "Public DNS name of the EC2 instance"
  value       = aws_instance.app.public_dns
}

output "rds_endpoint" {
  description = "RDS PostgreSQL endpoint"
  value       = aws_db_instance.postgres.endpoint
}

output "application_url" {
  description = "URL to access your application"
  value       = "http://${aws_instance.app.public_ip}:8080"
}

output "ssh_command" {
  description = "SSH command to connect to the EC2 instance"
  value       = "ssh -i ~/.ssh/cryple-key ec2-user@${aws_instance.app.public_ip}"
}
