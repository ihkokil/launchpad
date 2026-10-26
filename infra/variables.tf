variable "aws_region" {
  description = "AWS region to deploy infrastructure"
  type        = string
  default     = "us-east-1"
}

variable "environment" {
  description = "Deployment environment"
  type        = string
  default     = "production"
}

variable "cluster_name" {
  description = "EKS Cluster Name"
  type        = string
  default     = "launchpad-eks"
}

variable "github_org" {
  description = "GitHub username or organization for OIDC trust"
  type        = string
  default     = "ihkokil"
}

variable "vpc_cidr" {
  description = "VPC CIDR block"
  type        = string
  default     = "10.0.0.0/16"
}

