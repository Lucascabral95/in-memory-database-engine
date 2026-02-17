variable "aws_region" {
  description = "AWS region for all resources."
  type        = string
  default     = "us-east-1"
}

variable "project_name" {
  description = "Project name used in resource naming."
  type        = string
  default     = "in-memory-db"
}

variable "environment" {
  description = "Deployment environment."
  type        = string
  default     = "prod"
}

variable "instance_type" {
  description = "EC2 instance type for the API host."
  type        = string
  default     = "t3.small"
}

variable "database_url" {
  description = "External PostgreSQL DSN consumed by the API."
  type        = string
  sensitive   = true
}

variable "jwt_secret" {
  description = "JWT secret used by the API."
  type        = string
  sensitive   = true
}

variable "port" {
  description = "Public HTTP port exposed by the API."
  type        = string
  default     = "8080"
}

variable "redis_tcp_enabled" {
  description = "Enables custom Redis-like TCP server in the API."
  type        = bool
  default     = false
}

variable "redis_tcp_port" {
  description = "Port for Redis-like TCP server."
  type        = string
  default     = "6379"
}
