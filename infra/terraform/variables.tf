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

variable "db_max_open_conns" {
  description = "Database maximum number of open connections."
  type        = number
  default     = 20
}

variable "db_max_idle_conns" {
  description = "Database maximum number of idle connections."
  type        = number
  default     = 10
}

variable "db_conn_max_lifetime" {
  description = "Maximum amount of time a connection may be reused (Go duration format, e.g. 5m)."
  type        = string
  default     = "5m"
}

variable "db_conn_max_idle_time" {
  description = "Maximum amount of time a connection may be idle (Go duration format, e.g. 1m)."
  type        = string
  default     = "1m"
}

variable "prometheus_remote_write_url" {
  description = "Grafana Cloud Prometheus remote_write endpoint URL."
  type        = string
  default     = ""
}

variable "prometheus_remote_write_username" {
  description = "Grafana Cloud remote_write username (instance ID)."
  type        = string
  default     = ""
}

variable "prometheus_remote_write_password" {
  description = "Grafana Cloud remote_write API key (MetricsPublisher)."
  type        = string
  sensitive   = true
  default     = ""
}
