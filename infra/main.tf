# VPC Module
module "vpc" {
  source = "./modules/vpc"

  vpc_cidr           = var.vpc_cidr
  cluster_name       = var.cluster_name
  environment        = var.environment
  availability_zones = ["${var.aws_region}a", "${var.aws_region}b"]
}

# EKS Cluster Module
module "eks" {
  source = "./modules/eks"

  cluster_name    = var.cluster_name
  cluster_version = "1.30"
  vpc_id          = module.vpc.vpc_id
  subnet_ids      = module.vpc.private_subnet_ids
  instance_types  = ["t3.medium"]
  desired_size    = 2
  min_size        = 1
  max_size        = 4
  environment     = var.environment

  depends_on = [module.vpc]
}

# Keyless GitHub Actions OIDC Authentication Module
module "iam_oidc" {
  source = "./modules/iam_oidc"

  github_org  = var.github_org
  github_repo = "launchpad*"
  environment = var.environment
}
