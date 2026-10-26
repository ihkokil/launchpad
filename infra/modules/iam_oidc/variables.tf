variable "github_org" {
  description = "GitHub Organization or Username"
  type        = string
  default     = "ihkokil"
}

variable "github_repo" {
  description = "GitHub Repository filter (wildcard supported, e.g. 'launchpad*' or '*')"
  type        = string
  default     = "*"
}

variable "environment" {
  description = "Environment identifier"
  type        = string
  default     = "production"
}

