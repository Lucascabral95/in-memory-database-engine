data "aws_caller_identity" "current" {}

data "aws_vpc" "default" {
  default = true
}

data "aws_subnets" "default" {
  filter {
    name   = "vpc-id"
    values = [data.aws_vpc.default.id]
  }
}

locals {
  name_prefix          = "${var.project_name}-${var.environment}"
  ssm_parameter_prefix = "/${var.project_name}/${var.environment}"
  selected_subnet_id   = sort(data.aws_subnets.default.ids)[0]

  common_tags = {
    Project     = var.project_name
    Environment = var.environment
    ManagedBy   = "terraform"
  }
}
