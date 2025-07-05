# Cryple AWS Deployment

This Terraform configuration deploys the Cryple application to AWS using free tier resources.

## Architecture

- **EC2 t2.micro**: Runs your Docker containerized application
- **RDS db.t3.micro**: PostgreSQL database (free tier - 20GB storage)
- **VPC**: Isolated network environment
- **Security Groups**: Firewall rules for app and database

## Prerequisites

1. **AWS Account** with free tier eligibility
2. **AWS CLI** configured with your credentials
3. **Terraform** installed (>= 1.0)
4. **SSH key pair** for EC2 access

### Quick Installation (macOS with Homebrew)

```bash
# Install AWS CLI
brew install awscli

# Install Terraform
brew install terraform

# Verify installations
aws --version
terraform --version
```

## Setup Instructions

### 1. Configure AWS CLI

```bash
aws configure
# Enter your AWS Access Key ID, Secret Access Key, and region
```

### 2. Generate SSH Key (if you don't have one)

```bash
ssh-keygen -t rsa -b 4096 -f ~/.ssh/cryple-key
# This creates cryple-key (private) and cryple-key.pub (public)
```

### 3. Push Your Docker Image (Important!)

Before deploying, you need to push your Docker image to Docker Hub:

```bash
# Build and tag your image
docker build -t philippeberto/cryple-api:latest -f build/Dockerfile .

# Login to Docker Hub
docker login

# Push the image
docker push philippeberto/cryple-api:latest
```

**Note**: Update the image name in `user_data.sh` if you use a different repository.

### 4. Configure Terraform Variables

```bash
cd terraform
cp terraform.tfvars.example terraform.tfvars
```

Edit `terraform.tfvars` with your values:

```hcl
aws_region   = "us-west-2"
project_name = "cryple"
public_key   = "ssh-rsa AAAAB3NzaC1yc2E... your-email@example.com"  # Content of ~/.ssh/cryple-key.pub
db_password  = "YourSecurePassword123!"
```

### 5. Deploy Infrastructure

```bash
# Initialize Terraform
terraform init

# Plan the deployment
terraform plan

# Apply the configuration
terraform apply
```

Type `yes` when prompted.

### 6. Access Your Application

After deployment (takes ~5-10 minutes), you'll see outputs:

```
ec2_public_ip = "54.123.45.67"
application_url = "http://54.123.45.67:8080"
ssh_command = "ssh -i ~/.ssh/cryple-key ec2-user@54.123.45.67"
```

## Testing Your Deployment

### 1. Health Check

```bash
curl http://YOUR_EC2_IP:8080/users/check
# Should return HTTP 400 (expected for this endpoint without auth)
```

### 2. Create a User

```bash
curl -X POST http://YOUR_EC2_IP:8080/users \
  -H "Content-Type: application/json" \
  -d '{
    "user_address": "1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
    "password": "abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890"
  }'
```

### 3. SSH into EC2 (for debugging)

```bash
ssh -i ~/.ssh/cryple-key ec2-user@YOUR_EC2_IP

# Check application status
sudo systemctl status cryple
sudo docker-compose -f /opt/cryple/docker-compose.yml logs
```

## Monitoring and Troubleshooting

### Application Logs

```bash
# SSH into EC2
ssh -i ~/.ssh/cryple-key ec2-user@YOUR_EC2_IP

# Check Docker containers
sudo docker ps

# Check application logs
sudo docker-compose -f /opt/cryple/docker-compose.yml logs cryple

# Check deployment log
sudo cat /var/log/cryple-deploy.log

# Get logs from ec2
aws ec2 get-console-output --instance-id i-091bb446836f7dd89 --query 'Output' --output text | tail -50
```

### Health Check Script

```bash
# Run the health check
sudo /opt/cryple/health-check.sh
```

## Cost Optimization

This setup uses AWS Free Tier resources:

- **EC2 t2.micro**: 750 hours/month (always free for 1 year)
- **RDS db.t3.micro**: 750 hours/month (free for 1 year)
- **Storage**: 20GB (free tier limit)
- **Data Transfer**: 15GB out/month

**Estimated monthly cost after free tier**: ~$25-30

## Cleanup

To avoid charges, destroy resources when done testing:

```bash
terraform destroy
```

Type `yes` when prompted.

## Security Notes

- The security group allows HTTP (8080) access from anywhere
- SSH access is open to 0.0.0.0/0 (consider restricting to your IP)
- Database is in private subnets, only accessible from the application
- For production, consider:
  - Using HTTPS with SSL certificates
  - Restricting SSH access
  - Enabling database encryption
  - Using Parameter Store for secrets

## Troubleshooting

### Common Issues

1. **"Image not found"**: Make sure you pushed your Docker image to Docker Hub
2. **"Database connection failed"**: RDS takes 5-10 minutes to be ready
3. **"Permission denied"**: Check your SSH key path and permissions
4. **"Free tier exceeded"**: Verify you're using t2.micro and db.t3.micro instances

### Useful Commands

```bash
# Check Terraform state
terraform state list

# Get resource details
terraform state show aws_instance.app

# Force recreate EC2 instance
terraform taint aws_instance.app
terraform apply
```

## Advanced Debugging Commands

### EC2 Instance Management

```bash
# SSH into EC2 instance
ssh -i ~/.ssh/cryple-key ec2-user@YOUR_EC2_IP

# Check EC2 instance status and details
aws ec2 describe-instances --instance-ids i-YOUR_INSTANCE_ID

# Get EC2 console output for boot debugging
aws ec2 get-console-output --instance-id i-YOUR_INSTANCE_ID --query 'Output' --output text | tail -50

# Check system resources
ssh -i ~/.ssh/cryple-key ec2-user@YOUR_EC2_IP "free -h && df -h && top -bn1"
```

### Docker Container Debugging

```bash
# Check Docker containers status
ssh -i ~/.ssh/cryple-key ec2-user@YOUR_EC2_IP "sudo docker ps -a"

# View Docker Compose services
ssh -i ~/.ssh/cryple-key ec2-user@YOUR_EC2_IP "cd /opt/cryple && sudo docker-compose ps"

# View application logs (real-time)
ssh -i ~/.ssh/cryple-key ec2-user@YOUR_EC2_IP "cd /opt/cryple && sudo docker-compose logs -f cryple"

# View last 50 lines of logs
ssh -i ~/.ssh/cryple-key ec2-user@YOUR_EC2_IP "cd /opt/cryple && sudo docker-compose logs --tail=50 cryple"

# Restart application container
ssh -i ~/.ssh/cryple-key ec2-user@YOUR_EC2_IP "cd /opt/cryple && sudo docker-compose restart cryple"

# Rebuild and restart with latest image
ssh -i ~/.ssh/cryple-key ec2-user@YOUR_EC2_IP "cd /opt/cryple && sudo docker-compose pull && sudo docker-compose down && sudo docker-compose up -d"

# Execute commands inside running container
ssh -i ~/.ssh/cryple-key ec2-user@YOUR_EC2_IP "sudo docker exec -it cryple-cryple-1 /bin/sh"

# Check container file system
ssh -i ~/.ssh/cryple-key ec2-user@YOUR_EC2_IP "sudo docker exec cryple-cryple-1 ls -la /app"
ssh -i ~/.ssh/cryple-key ec2-user@YOUR_EC2_IP "sudo docker exec cryple-cryple-1 ls -la /app/migrations"
```

### Database Debugging

```bash
# Test database connectivity from EC2
ssh -i ~/.ssh/cryple-key ec2-user@YOUR_EC2_IP "sudo docker exec cryple-cryple-1 nc -zv YOUR_RDS_ENDPOINT 5432"

# Get RDS instance details
aws rds describe-db-instances --db-instance-identifier cryple-db

# Check RDS parameter group (for SSL settings)
aws rds describe-db-parameter-groups --db-parameter-group-name cryple-db-params

# Reboot RDS instance (if needed for parameter changes)
aws rds reboot-db-instance --db-instance-identifier cryple-db

# Connect to database directly (if needed)
ssh -i ~/.ssh/cryple-key ec2-user@YOUR_EC2_IP "sudo docker run --rm -it postgres:15 psql -h YOUR_RDS_ENDPOINT -U YOUR_DB_USER -d cryple"
```

### Application Testing

```bash
# Test application health
curl -v http://YOUR_EC2_IP:8080/users/check

# Test with detailed output
curl -i -X GET http://YOUR_EC2_IP:8080/users/check

# Test user creation
curl -X POST http://YOUR_EC2_IP:8080/users \
  -H "Content-Type: application/json" \
  -d '{
    "user_address": "1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
    "password": "abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890"
  }'

# Test from within EC2 (useful for network debugging)
ssh -i ~/.ssh/cryple-key ec2-user@YOUR_EC2_IP "curl -v http://localhost:8080/users/check"
```

### Configuration Management

```bash
# Check current docker-compose.yml
ssh -i ~/.ssh/cryple-key ec2-user@YOUR_EC2_IP "cat /opt/cryple/docker-compose.yml"

# View deployment logs
ssh -i ~/.ssh/cryple-key ec2-user@YOUR_EC2_IP "sudo cat /var/log/cryple-deploy.log"

# Update environment variables (example)
ssh -i ~/.ssh/cryple-key ec2-user@YOUR_EC2_IP 'cd /opt/cryple && sudo tee docker-compose.yml << EOF
version: "3.9"
services:
  cryple:
    image: philippeberto/cryple-api:0.0.1
    ports:
      - "8080:8080"
    environment:
      - POSTGRES_MIGRATION=true
      - DEBUG=true
      # ... other environment variables
    restart: unless-stopped
EOF'
```

### Local Development & Testing

```bash
# Build and test Docker image locally
docker build -t cryple-test -f build/Dockerfile .
docker run --rm cryple-test ls -la /app/migrations

# Test local Docker image with environment variables
docker run --rm -e POSTGRES_MIGRATION=false cryple-test

# Push updated image to registry
docker tag cryple-test:latest philippeberto/cryple-api:0.0.1
docker push philippeberto/cryple-api:0.0.1
```

### Log Analysis

```bash
# Search for specific errors in logs
ssh -i ~/.ssh/cryple-key ec2-user@YOUR_EC2_IP "cd /opt/cryple && sudo docker-compose logs cryple | grep -i error"

# Search for migration-related logs
ssh -i ~/.ssh/cryple-key ec2-user@YOUR_EC2_IP "cd /opt/cryple && sudo docker-compose logs cryple | grep -i migration"

# Search for database connection logs
ssh -i ~/.ssh/cryple-key ec2-user@YOUR_EC2_IP "cd /opt/cryple && sudo docker-compose logs cryple | grep -i 'database\|postgres'"

# Monitor logs in real-time with filtering
ssh -i ~/.ssh/cryple-key ec2-user@YOUR_EC2_IP "cd /opt/cryple && sudo docker-compose logs -f cryple | grep -E '(ERROR|WARN|migration|database)'"
```

### Performance Monitoring

```bash
# Check container resource usage
ssh -i ~/.ssh/cryple-key ec2-user@YOUR_EC2_IP "sudo docker stats cryple-cryple-1 --no-stream"

# Check EC2 instance metrics
ssh -i ~/.ssh/cryple-key ec2-user@YOUR_EC2_IP "top -bn1 | head -20"
ssh -i ~/.ssh/cryple-key ec2-user@YOUR_EC2_IP "iostat -x 1 5"

# Check application response times
time curl -s http://YOUR_EC2_IP:8080/users/check
```

### Quick Deployment Fixes

```bash
# Quick restart after code changes
ssh -i ~/.ssh/cryple-key ec2-user@YOUR_EC2_IP "cd /opt/cryple && sudo docker-compose down && sudo docker-compose up -d"

# Force pull latest image and restart
ssh -i ~/.ssh/cryple-key ec2-user@YOUR_EC2_IP "cd /opt/cryple && sudo docker-compose pull && sudo docker-compose down && sudo docker-compose up -d"

# Emergency stop
ssh -i ~/.ssh/cryple-key ec2-user@YOUR_EC2_IP "cd /opt/cryple && sudo docker-compose down"

# Check if service is responding
ssh -i ~/.ssh/cryple-key ec2-user@YOUR_EC2_IP "timeout 5 curl -s http://localhost:8080/users/check && echo 'Service is responding' || echo 'Service is not responding'"
```
