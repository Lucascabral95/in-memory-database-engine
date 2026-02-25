resource "aws_ssm_parameter" "database_url" {
  name      = "${local.ssm_parameter_prefix}/DATABASE_URL"
  type      = "SecureString"
  value     = var.database_url
  overwrite = true
}

resource "aws_ssm_parameter" "jwt_secret" {
  name      = "${local.ssm_parameter_prefix}/JWT_SECRET"
  type      = "SecureString"
  value     = var.jwt_secret
  overwrite = true
}

resource "aws_ssm_parameter" "port" {
  name      = "${local.ssm_parameter_prefix}/PORT"
  type      = "String"
  value     = var.port
  overwrite = true
}

resource "aws_ssm_parameter" "environment" {
  name      = "${local.ssm_parameter_prefix}/ENV"
  type      = "String"
  value     = var.environment
  overwrite = true
}

resource "aws_ssm_parameter" "redis_tcp_enabled" {
  name      = "${local.ssm_parameter_prefix}/REDIS_TCP_ENABLED"
  type      = "String"
  value     = tostring(var.redis_tcp_enabled)
  overwrite = true
}

resource "aws_ssm_parameter" "redis_tcp_port" {
  name      = "${local.ssm_parameter_prefix}/REDIS_TCP_PORT"
  type      = "String"
  value     = var.redis_tcp_port
  overwrite = true
}

resource "aws_ssm_parameter" "db_max_open_conns" {
  name      = "${local.ssm_parameter_prefix}/DB_MAX_OPEN_CONNS"
  type      = "String"
  value     = tostring(var.db_max_open_conns)
  overwrite = true
}

resource "aws_ssm_parameter" "db_max_idle_conns" {
  name      = "${local.ssm_parameter_prefix}/DB_MAX_IDLE_CONNS"
  type      = "String"
  value     = tostring(var.db_max_idle_conns)
  overwrite = true
}

resource "aws_ssm_parameter" "db_conn_max_lifetime" {
  name      = "${local.ssm_parameter_prefix}/DB_CONN_MAX_LIFETIME"
  type      = "String"
  value     = var.db_conn_max_lifetime
  overwrite = true
}

resource "aws_ssm_parameter" "db_conn_max_idle_time" {
  name      = "${local.ssm_parameter_prefix}/DB_CONN_MAX_IDLE_TIME"
  type      = "String"
  value     = var.db_conn_max_idle_time
  overwrite = true
}

resource "aws_ssm_parameter" "prometheus_remote_write_url" {
  count     = var.prometheus_remote_write_url != "" ? 1 : 0
  name      = "${local.ssm_parameter_prefix}/PROMETHEUS_REMOTE_WRITE_URL"
  type      = "String"
  value     = var.prometheus_remote_write_url
  overwrite = true
}

resource "aws_ssm_parameter" "prometheus_remote_write_username" {
  count     = var.prometheus_remote_write_username != "" ? 1 : 0
  name      = "${local.ssm_parameter_prefix}/PROMETHEUS_REMOTE_WRITE_USERNAME"
  type      = "String"
  value     = var.prometheus_remote_write_username
  overwrite = true
}

resource "aws_ssm_parameter" "prometheus_remote_write_password" {
  count     = var.prometheus_remote_write_password != "" ? 1 : 0
  name      = "${local.ssm_parameter_prefix}/PROMETHEUS_REMOTE_WRITE_PASSWORD"
  type      = "SecureString"
  value     = var.prometheus_remote_write_password
  overwrite = true
}
