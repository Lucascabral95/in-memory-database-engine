data "aws_iam_policy_document" "ec2_assume_role" {
  statement {
    effect = "Allow"

    principals {
      type        = "Service"
      identifiers = ["ec2.amazonaws.com"]
    }

    actions = ["sts:AssumeRole"]
  }
}

resource "aws_iam_role" "ec2_role" {
  name               = "${local.name_prefix}-ec2-role"
  assume_role_policy = data.aws_iam_policy_document.ec2_assume_role.json
}

resource "aws_iam_role_policy_attachment" "ssm_managed_core" {
  role       = aws_iam_role.ec2_role.name
  policy_arn = "arn:aws:iam::aws:policy/AmazonSSMManagedInstanceCore"
}

resource "aws_iam_role_policy_attachment" "ecr_readonly" {
  role       = aws_iam_role.ec2_role.name
  policy_arn = "arn:aws:iam::aws:policy/AmazonEC2ContainerRegistryReadOnly"
}

data "aws_iam_policy_document" "ec2_ssm_access" {
  statement {
    sid    = "ReadProjectParameters"
    effect = "Allow"
    actions = [
      "ssm:GetParameter",
      "ssm:GetParameters"
    ]
    resources = [
      "arn:aws:ssm:${var.aws_region}:${data.aws_caller_identity.current.account_id}:parameter${local.ssm_parameter_prefix}/*"
    ]
  }

  statement {
    sid    = "ReadProjectParametersByPath"
    effect = "Allow"
    actions = [
      "ssm:GetParametersByPath"
    ]
    resources = [
      "arn:aws:ssm:${var.aws_region}:${data.aws_caller_identity.current.account_id}:parameter${local.ssm_parameter_prefix}",
      "arn:aws:ssm:${var.aws_region}:${data.aws_caller_identity.current.account_id}:parameter${local.ssm_parameter_prefix}/*"
    ]
  }

  statement {
    sid       = "DecryptSecureString"
    effect    = "Allow"
    actions   = ["kms:Decrypt"]
    resources = ["*"]
  }
}

resource "aws_iam_policy" "ec2_ssm_access" {
  name   = "${local.name_prefix}-ec2-ssm-access"
  policy = data.aws_iam_policy_document.ec2_ssm_access.json
}

resource "aws_iam_role_policy_attachment" "ec2_ssm_access" {
  role       = aws_iam_role.ec2_role.name
  policy_arn = aws_iam_policy.ec2_ssm_access.arn
}

resource "aws_iam_instance_profile" "ec2_profile" {
  name = "${local.name_prefix}-instance-profile"
  role = aws_iam_role.ec2_role.name
}
