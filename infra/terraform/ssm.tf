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
