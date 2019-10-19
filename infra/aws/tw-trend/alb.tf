resource "aws_security_group_rule" "instance_in_alb" {
  type                     = "ingress"
  from_port                = 32768
  to_port                  = 61000
  protocol                 = "tcp"
  source_security_group_id = module.alb_sg_https.this_security_group_id
  security_group_id        = aws_security_group.instance.id
}

module "alb_sg_https" {
  source  = "terraform-aws-modules/security-group/aws"
  version = "3.1.0"

  name   = "${var.name}-${var.env}-alb"
  vpc_id = module.vpc.vpc_id

  ingress_with_cidr_blocks = [
    {
      rule        = "https-443-tcp"
      cidr_blocks = "0.0.0.0/0"
    },
    {
      rule        = "http-80-tcp"
      cidr_blocks = "0.0.0.0/0"
    },

  ]

  egress_with_cidr_blocks = [
    {
      rule        = "all-tcp"
      cidr_blocks = "0.0.0.0/0"
    },
  ]

  tags = "${var.tags}"
}

module "alb" {
  source  = "terraform-aws-modules/alb/aws"
  version = "~> v4.0"

  load_balancer_name            = "${var.name}-${var.env}-alb"
  security_groups = [module.alb_sg_https.this_security_group_id]

  http_tcp_listeners            = "${list(map("port", "80", "protocol", "HTTP"))}"
  http_tcp_listeners_count      = "1"

  target_groups                 = "${list(map("name", "tw-trends", "backend_protocol", "HTTP", "backend_port", "5000"))}"
  target_groups_count           = "1"

  logging_enabled          = true
  log_bucket_name          = aws_s3_bucket.log_bucket.id
  log_location_prefix      = "alb"

  subnets = module.vpc.public_subnets

  tags = var.tags
  vpc_id = module.vpc.vpc_id
}

data "aws_iam_policy_document" "bucket_policy_logs_alb" {
  statement {
    sid       = "AllowToPutLoadBalancerLogsToS3Bucket"
    actions   = ["s3:PutObject"]
    resources = ["arn:aws:s3:::${var.name}-${var.env}-alb-logs/alb/AWSLogs/${data.aws_caller_identity.current.account_id}/*"]

    principals {
      type        = "AWS"
      identifiers = ["arn:aws:iam::${data.aws_elb_service_account.main.id}:root"]
    }
  }
}

resource "aws_s3_bucket" "log_bucket" {
  bucket        = "${var.name}-${var.env}-alb-logs"
  policy        = data.aws_iam_policy_document.bucket_policy_logs_alb.json
  force_destroy = true
  tags          = var.tags

  lifecycle_rule {
    id      = "log-expiration"
    enabled = "true"

    expiration {
      days = "7"
    }
  }
}