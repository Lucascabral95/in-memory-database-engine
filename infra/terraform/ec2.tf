data "aws_ssm_parameter" "amazon_linux_2023_ami" {
  name = "/aws/service/ami-amazon-linux-latest/al2023-ami-kernel-default-x86_64"
}

resource "aws_security_group" "public_open" {
  name        = "${local.name_prefix}-public-open-sg"
  description = "Open to all traffic by explicit project requirement."
  vpc_id      = data.aws_vpc.default.id

  ingress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

resource "aws_instance" "api" {
  ami                         = data.aws_ssm_parameter.amazon_linux_2023_ami.value
  instance_type               = var.instance_type
  subnet_id                   = local.selected_subnet_id
  vpc_security_group_ids      = [aws_security_group.public_open.id]
  iam_instance_profile        = aws_iam_instance_profile.ec2_profile.name
  associate_public_ip_address = true
  user_data_replace_on_change = true

  user_data = templatefile("${path.module}/user_data.sh.tftpl", {
    aws_region                = var.aws_region
    ssm_parameter_prefix      = local.ssm_parameter_prefix
    ecr_repository_url        = aws_ecr_repository.api.repository_url
    ecr_registry              = split("/", aws_ecr_repository.api.repository_url)[0]
    app_port                  = var.port
    redis_port                = var.redis_tcp_port
    container_name            = "${local.name_prefix}-api"
    prometheus_container_name = "${local.name_prefix}-prometheus"
    prometheus_image          = "prom/prometheus"
  })

  root_block_device {
    volume_type = "gp3"
    volume_size = 20
  }

  metadata_options {
    http_endpoint = "enabled"
    http_tokens   = "required"
  }

  tags = merge(local.common_tags, {
    Name = "${local.name_prefix}-ec2"
  })
}
