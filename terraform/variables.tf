variable "aws_region" {
  description = "AWS region"
  type        = string
  default     = "us-west-2"

  validation {
    condition = contains([
      "us-west-2",      # Oregon (recommended)
      "us-east-1",      # N. Virginia 
      "us-west-1",      # N. California
      "eu-west-1",      # Ireland
      "eu-central-1",   # Frankfurt
      "ap-southeast-1", # Singapore
      "ap-northeast-1", # Tokyo
    ], var.aws_region)
    error_message = "Region must be one of the supported regions."
  }
}

variable "project_name" {
  description = "Name of the project (used for resource naming)"
  type        = string
  default     = "cryple"
}

variable "public_key" {
  description = "Public key for EC2 instance access"
  type        = string
}

variable "db_password" {
  description = "Password for the RDS PostgreSQL instance"
  type        = string
  sensitive   = true
}

variable "ami_id" {
  description = "AMI ID for the EC2 instance (Amazon Linux 2023)"
  type        = string
  default     = ""
}

# AMI mapping for Amazon Linux 2023 by region
locals {
  ami_ids = {
    us-west-2      = "ami-0aff18ec83b712f05" # Oregon
    us-east-1      = "ami-0b72821e2f351e396" # N. Virginia
    us-west-1      = "ami-0827b6c5b977c020e" # N. California
    eu-west-1      = "ami-0c02fb55956c7d316" # Ireland
    eu-central-1   = "ami-04f76ebf53292ef4d" # Frankfurt
    ap-southeast-1 = "ami-0497a974f8d5dcef8" # Singapore
    ap-northeast-1 = "ami-0d52744d6551d851e" # Tokyo
  }

  selected_ami = var.ami_id != "" ? var.ami_id : local.ami_ids[var.aws_region]
}
