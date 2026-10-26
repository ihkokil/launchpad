terraform {
  required_version = ">= 1.5.0"

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.50"
    }
  }

  # For staging/production, configure remote S3 backend with DynamoDB locking
  # backend "s3" {
  #   bucket         = "launchpad-terraform-state-ihkokil"
  #   key            = "platform/terraform.tfstate"
  #   region         = "us-east-1"
  #   dynamodb_table = "launchpad-terraform-locks"
  #   encrypt        = true
  # }
}

provider "aws" {
  region = var.aws_region

  default_tags {
    tags = {
      Project     = "Launchpad-IDP"
      ManagedBy   = "Terraform"
      Owner       = var.github_org
      Environment = var.environment
    }
  }
}
