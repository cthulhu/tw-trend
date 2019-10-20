resource "aws_key_pair" "ecs" {
  key_name   = "${var.name}-${var.env}-ecs"
  public_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQDh+LprE1iivuba6N/Vx/XyAZCFBSH4bi3J7vwyxMU61QrVEYnlb+nycUgILMzrHLsQ1f3qboFdNv3RfFNQULiWI99hnEpnaOqBByoUaW1Ejq/fVg+7Rsg3QpZRQQPBimDcNTayFG5OFvK8FgeY5oZbbAsqNlxR3PgwmMhBPz94Zv4TfDUSBV+0490MaSkjuC/o3BDP/PQykjb5kZ0KtuJyyQ67xlbXBS44NUefXH4Ym9NN2mU4BncKZGQO6MXRDnKHOx6sMkwHFo1YcLIcQzi0Ukmf5Pni9Lj6IdCYtC3ntjQaqCEj4oZNxGYNx3VaspbhCkzT2aPlnlxgsPa0MxnzCSh//Eh+xIZpJ9DGV5HJbehXe5WTdm96kJG8l4yUuQvL7101A1lUjCxCXINmqVBPPEtJCbeaeoclV+aGTSoiml6umdAJe7wHAOaI3zUigJ+A3BfxdloDsIt5d1h7IbR2PrhEdXGCZFY9eOwaGt47fbF6u5p6LWatKa8x/9ZBk54plKPkiB4G6Jdt5uUSkLwX2b5meXaMxQeWLMrf3X/zCCdqWEkZrx+bldS8Kx7O/Wg8DNuBzx+84SBfpoEqNl7F9TJ/pdsEw1S+baQfWMz4Km+26uQHwCFxC+24w96O0HT9Y56Lvb4Izu5QDRGaMcLL73sJlBPNKf0M02JEEAwwIw== spogrebnyak@xebia.com"
}

resource "aws_ecs_cluster" "ecs_cluster" {
  name = "${var.name}-${var.env}-cluster"
  tags = var.tags
}

data "template_file" "task_definition" {
  template = file("${path.module}/ecs/task-definition.json")

  vars = {
    image_url        = "stanpogrebnyak/tw-trend:latest"
    container_name   = "tw-trend-container"
    log_group_region = var.region
    log_group_name   = aws_cloudwatch_log_group.app.name
    container_port   = var.container_port
    name             = var.name
    env              = var.env
    account_id       = data.aws_caller_identity.current.account_id
    region           = var.region
  }
}

resource "aws_ecs_task_definition" "ecs_task_definition" {
  family                = "${var.name}-${var.env}"
  container_definitions = "${data.template_file.task_definition.rendered}"
  tags                  = var.tags
  execution_role_arn    = aws_iam_role.ecs_task.arn
  task_role_arn         = aws_iam_role.ecs_task.arn
  volume {
    name      = "data-dir"
    host_path = "/mnt/data"
  }
}


data "aws_iam_policy_document" "instance_policy" {
  statement {
    sid = "CloudwatchPutMetricData"

    actions = [
      "cloudwatch:PutMetricData",
    ]

    resources = [
      "*",
    ]
  }

  statement {
    sid = "InstanceLogging"

    actions = [
      "logs:CreateLogGroup",
      "logs:CreateLogStream",
      "logs:PutLogEvents",
      "logs:DescribeLogStreams",
    ]
    resources = [aws_cloudwatch_log_group.instance.arn]
  }
}
resource "aws_iam_policy" "instance_policy" {
  name   = "${var.name}-${var.env}-ecs-instance"
  path   = "/"
  policy = "${data.aws_iam_policy_document.instance_policy.json}"
}

resource "aws_iam_role" "instance" {
  name = "${var.name}-${var.env}-instance-role"

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Principal": {
        "Service": "ec2.amazonaws.com"
      },
      "Effect": "Allow",
      "Sid": ""
    }
  ]
}
EOF
}

resource "aws_iam_role_policy_attachment" "ecs_policy" {
  role       = "${aws_iam_role.instance.name}"
  policy_arn = "arn:aws:iam::aws:policy/service-role/AmazonEC2ContainerServiceforEC2Role"
}

resource "aws_iam_role_policy_attachment" "instance_policy" {
  role       = "${aws_iam_role.instance.name}"
  policy_arn = "${aws_iam_policy.instance_policy.arn}"
}

resource "aws_iam_instance_profile" "instance" {
  name = "${var.name}-${var.env}-instance-profile"
  role = "${aws_iam_role.instance.name}"
}

resource "aws_security_group" "instance" {
  name        = "${var.name}-${var.env}-container-instance"
  description = "Security Group managed by Terraform"
  vpc_id      = "${module.vpc.vpc_id}"
}

resource "aws_security_group_rule" "instance_out_all" {
  type              = "egress"
  from_port         = 0
  to_port           = 65535
  protocol          = "tcp"
  cidr_blocks       = ["0.0.0.0/0"]
  security_group_id = "${aws_security_group.instance.id}"
}

resource "aws_security_group_rule" "allow_ssh" {
  type            = "ingress"
  from_port       = 22
  to_port         = 22
  protocol        = "tcp"
  security_group_id = aws_security_group.instance.id
  source_security_group_id = aws_security_group.bastion.id
}

data "template_file" "user_data" {
  template = file("${path.module}/ecs/user_data.sh")

  vars = {
    ecs_cluster = aws_ecs_cluster.ecs_cluster.name
    log_group   = aws_cloudwatch_log_group.instance.name
  }
}

data "aws_ami" "ecs" {
  most_recent = true

  filter {
    name   = "name"
    values = ["*amazon-ecs-optimized"]
  }

  filter {
    name   = "virtualization-type"
    values = ["hvm"]
  }

  owners = ["amazon"]
}

resource "aws_launch_configuration" "instance" {
  name_prefix          = "${var.name}-${var.env}-launch-configuration"
  image_id             = data.aws_ami.ecs.id
  instance_type        = var.instance_type
  iam_instance_profile = aws_iam_instance_profile.instance.name
  user_data            = data.template_file.user_data.rendered
  security_groups      = [aws_security_group.instance.id]
  key_name             = aws_key_pair.ecs.key_name

  root_block_device {
    volume_size = var.instance_root_volume_size
    volume_type = "gp2"
  }

  ebs_block_device {
    encrypted   = true
    device_name = "/dev/xvdf"
    volume_size = var.instance_root_volume_size
  }

  lifecycle {
    create_before_destroy = true
  }
}

resource "aws_autoscaling_group" "asg" {
  name = "${var.name}-${var.env}-asg"

  launch_configuration = "${aws_launch_configuration.instance.name}"
  vpc_zone_identifier  = module.vpc.private_subnets
  max_size             = "${var.asg_max_size}"
  min_size             = "${var.asg_min_size}"
  desired_capacity     = "${var.asg_desired_size}"

  health_check_grace_period = 300
  health_check_type         = "EC2"

  lifecycle {
    create_before_destroy = true
  }
}


resource "aws_iam_role" "ecs_task" {
  name = "${var.name}-${var.env}-ecs-task"
  path = "/ecs/"

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "Service": [
          "ecs-tasks.amazonaws.com"
        ]
      },
      "Action": "sts:AssumeRole"
    }
  ]
}
EOF

}

resource "aws_iam_policy" "ecs_task" {
  name = "${var.name}-${var.env}-ecs-task"
  path = "/"

  policy = <<EOF
{
    "Version": "2012-10-17",
    "Statement": [
    {
      "Action": [
        "ssm:DescribeParameters",
        "kms:Decrypt",
        "secretsmanager:GetSecretValue"
      ],
      "Effect": "Allow",
      "Resource": "*"
    },
    {
      "Effect": "Allow",
      "Action": [
        "logs:CreateLogStream",
        "logs:PutLogEvents"
      ],
      "Resource": "*"
    },
    {
      "Action": [
        "ssm:GetParameters"
      ],
      "Effect": "Allow",
      "Resource": "arn:aws:ssm:${var.region}:${data.aws_caller_identity.current.account_id}:parameter/${var.name}/${var.env}/*"
    },
    {
      "Effect": "Allow",
      "Action": "dynamodb:*",
      "Resource": "*"
    }
  ]
}
EOF

}

resource "aws_iam_policy_attachment" "ecs_task" {
  name       = "${var.name}-${var.env}-ecs-task"
  roles      = [aws_iam_role.ecs_task.name]
  policy_arn = aws_iam_policy.ecs_task.arn
}


