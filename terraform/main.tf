# Configure the AWS Provider
terraform {
  required_version = ">= 1.0"
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}

provider "aws" {
  region = var.aws_region
}

# Data source for availability zones
data "aws_availability_zones" "available" {
  state = "available"
}

# VPC
resource "aws_vpc" "main" {
  cidr_block           = "10.0.0.0/16"
  enable_dns_hostnames = true
  enable_dns_support   = true

  tags = {
    Name = "${var.project_name}-vpc"
  }
}

# Internet Gateway
resource "aws_internet_gateway" "main" {
  vpc_id = aws_vpc.main.id

  tags = {
    Name = "${var.project_name}-igw"
  }
}

# Public Subnet
resource "aws_subnet" "public" {
  vpc_id                  = aws_vpc.main.id
  cidr_block              = "10.0.1.0/24"
  availability_zone       = data.aws_availability_zones.available.names[0]
  map_public_ip_on_launch = true

  tags = {
    Name = "${var.project_name}-public-subnet"
  }
}

# Private Subnet 1
resource "aws_subnet" "private_1" {
  vpc_id            = aws_vpc.main.id
  cidr_block        = "10.0.2.0/24"
  availability_zone = data.aws_availability_zones.available.names[0]

  tags = {
    Name = "${var.project_name}-private-subnet-1"
  }
}

# Private Subnet 2 (required for RDS subnet group)
resource "aws_subnet" "private_2" {
  vpc_id            = aws_vpc.main.id
  cidr_block        = "10.0.3.0/24"
  availability_zone = data.aws_availability_zones.available.names[1]

  tags = {
    Name = "${var.project_name}-private-subnet-2"
  }
}

# Route Table for Public Subnet
resource "aws_route_table" "public" {
  vpc_id = aws_vpc.main.id

  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.main.id
  }

  tags = {
    Name = "${var.project_name}-public-rt"
  }
}

# Route Table Association for Public Subnet
resource "aws_route_table_association" "public" {
  subnet_id      = aws_subnet.public.id
  route_table_id = aws_route_table.public.id
}

# Security Group for EC2 (Application)
resource "aws_security_group" "app" {
  name_prefix = "${var.project_name}-app-"
  vpc_id      = aws_vpc.main.id

  # HTTP access
  ingress {
    from_port   = 8080
    to_port     = 8080
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  # SSH access
  ingress {
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"] # Consider restricting this to your IP
  }

  # All outbound traffic
  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name = "${var.project_name}-app-sg"
  }
}

# Security Group for RDS
resource "aws_security_group" "rds" {
  name_prefix = "${var.project_name}-rds-"
  vpc_id      = aws_vpc.main.id

  # PostgreSQL access from app security group
  ingress {
    from_port       = 5432
    to_port         = 5432
    protocol        = "tcp"
    security_groups = [aws_security_group.app.id]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name = "${var.project_name}-rds-sg"
  }
}

# Key Pair for EC2 access
resource "aws_key_pair" "main" {
  key_name   = "${var.project_name}-key"
  public_key = var.public_key
}

# EC2 Instance
resource "aws_instance" "app" {
  ami                    = local.selected_ami
  instance_type          = "t2.micro" # Free tier eligible
  key_name               = aws_key_pair.main.key_name
  vpc_security_group_ids = [aws_security_group.app.id]
  subnet_id              = aws_subnet.public.id

  user_data = base64encode(templatefile("${path.module}/user_data.sh", {
    db_host     = aws_db_instance.postgres.address
    db_port     = aws_db_instance.postgres.port
    db_name     = aws_db_instance.postgres.db_name
    db_user     = aws_db_instance.postgres.username
    db_password = var.db_password
  }))

  tags = {
    Name = "${var.project_name}-app"
  }
}

# RDS Subnet Group
resource "aws_db_subnet_group" "main" {
  name       = "${var.project_name}-db-subnet-group"
  subnet_ids = [aws_subnet.private_1.id, aws_subnet.private_2.id]

  tags = {
    Name = "${var.project_name}-db-subnet-group"
  }
}

# RDS Parameter Group to disable SSL requirement
resource "aws_db_parameter_group" "postgres" {
  family = "postgres17"
  name   = "${var.project_name}-postgres-params"

  parameter {
    name  = "rds.force_ssl"
    value = "0"
  }

  tags = {
    Name = "${var.project_name}-postgres-params"
  }
}

# RDS PostgreSQL Instance
resource "aws_db_instance" "postgres" {
  identifier     = "${var.project_name}-db"
  engine         = "postgres"
  engine_version = "17.5"
  instance_class = "db.t3.micro" # Free tier eligible

  allocated_storage     = 20 # Free tier: up to 20GB
  max_allocated_storage = 20
  storage_encrypted     = false # Free tier doesn't support encryption

  db_name  = "cryple"
  username = "FCIVpYdclzpDiL0K"
  password = var.db_password

  vpc_security_group_ids = [aws_security_group.rds.id]
  db_subnet_group_name   = aws_db_subnet_group.main.name

  skip_final_snapshot = true  # For testing purposes
  deletion_protection = false # For testing purposes

  backup_retention_period = 0 # Disable backups for free tier

  # SSL Configuration
  parameter_group_name = aws_db_parameter_group.postgres.name

  tags = {
    Name = "${var.project_name}-postgres"
  }
}
