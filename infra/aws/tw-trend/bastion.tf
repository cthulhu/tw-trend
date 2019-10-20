resource "aws_key_pair" "bastion" {
  key_name   = "${var.name}-${var.env}-bastion"
  public_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQDh+LprE1iivuba6N/Vx/XyAZCFBSH4bi3J7vwyxMU61QrVEYnlb+nycUgILMzrHLsQ1f3qboFdNv3RfFNQULiWI99hnEpnaOqBByoUaW1Ejq/fVg+7Rsg3QpZRQQPBimDcNTayFG5OFvK8FgeY5oZbbAsqNlxR3PgwmMhBPz94Zv4TfDUSBV+0490MaSkjuC/o3BDP/PQykjb5kZ0KtuJyyQ67xlbXBS44NUefXH4Ym9NN2mU4BncKZGQO6MXRDnKHOx6sMkwHFo1YcLIcQzi0Ukmf5Pni9Lj6IdCYtC3ntjQaqCEj4oZNxGYNx3VaspbhCkzT2aPlnlxgsPa0MxnzCSh//Eh+xIZpJ9DGV5HJbehXe5WTdm96kJG8l4yUuQvL7101A1lUjCxCXINmqVBPPEtJCbeaeoclV+aGTSoiml6umdAJe7wHAOaI3zUigJ+A3BfxdloDsIt5d1h7IbR2PrhEdXGCZFY9eOwaGt47fbF6u5p6LWatKa8x/9ZBk54plKPkiB4G6Jdt5uUSkLwX2b5meXaMxQeWLMrf3X/zCCdqWEkZrx+bldS8Kx7O/Wg8DNuBzx+84SBfpoEqNl7F9TJ/pdsEw1S+baQfWMz4Km+26uQHwCFxC+24w96O0HT9Y56Lvb4Izu5QDRGaMcLL73sJlBPNKf0M02JEEAwwIw== spogrebnyak@xebia.com"
}

resource "aws_security_group" "bastion" {
  name        = "${var.name}-${var.env}-bastion"
  vpc_id      = module.vpc.vpc_id
  description = "Bastion security group (only SSH inbound access is allowed)"

  tags = var.tags
}

resource "aws_security_group_rule" "ssh_ingress" {
  type              = "ingress"
  from_port         = "22"
  to_port           = "22"
  protocol          = "tcp"
  cidr_blocks       = ["0.0.0.0/0"]
  ipv6_cidr_blocks  = ["::/0"]
  security_group_id = aws_security_group.bastion.id
}

# resource "aws_security_group_rule" "ssh_sg_ingress" {
#   type                     = "ingress"
#   from_port                = "22"
#   to_port                  = "22"
#   protocol                 = "tcp"
#   source_security_group_id = "${element(var.allowed_security_groups, count.index)}"
#   security_group_id        = "${aws_security_group.bastion.id}"
# }

resource "aws_security_group_rule" "bastion_all_egress" {
  type      = "egress"
  from_port = "0"
  to_port   = "65535"
  protocol  = "all"

  cidr_blocks = [
    "0.0.0.0/0",
  ]

  ipv6_cidr_blocks = [
    "::/0",
  ]

  security_group_id = aws_security_group.bastion.id
}

data "template_file" "bastion_user_data" {
  template = "${file("${path.module}/bastion/user_data.sh")}"
  vars = {
    s3_bucket_name              = ""
    s3_bucket_uri               = ""
    ssh_user                    = "ec2-user"
    keys_update_frequency       = "*/5 * * * *"
    enable_hourly_cron_updates  = "true"
    additional_user_data_script = ""
  }
}

resource "aws_iam_instance_profile" "bastion" {
  name = "${var.name}-${var.env}-bastion"
  role = "${aws_iam_role.bastion.name}"
}

resource "aws_iam_role" "bastion" {
  name = "${var.name}-${var.env}-bastion"
  path = "/"

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Sid": "",
      "Effect": "Allow",
      "Principal": {
        "Service": "ec2.amazonaws.com"
      },
      "Action": "sts:AssumeRole"
    }
  ]
}
EOF
}

resource "aws_iam_role_policy" "bastion" {
  name = "${var.name}-${var.env}-bastion"
  role = "${aws_iam_role.bastion.id}"

  policy = <<EOF
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Sid": "Stmt1425916919000",
            "Effect": "Allow",
            "Action": [
                "s3:List*",
                "s3:Get*"
            ],
            "Resource": "*"
        }
    ]
}
EOF
}

resource "aws_launch_configuration" "bastion" {
  name_prefix       = "${var.name}-${var.env}-bastion"
  image_id          = "ami-07683a44e80cd32c5"
  instance_type     = "t2.micro"
  user_data         = data.template_file.bastion_user_data.rendered

  security_groups = [
    aws_security_group.bastion.id,
    aws_security_group.instance.id,
  ]

  root_block_device {
    volume_size = "20"
  }

  iam_instance_profile        = aws_iam_instance_profile.bastion.name
  associate_public_ip_address = true
  key_name                    = aws_key_pair.bastion.key_name

  lifecycle {
    create_before_destroy = true
  }
}

resource "aws_autoscaling_group" "bastion" {
  name = "${var.name}-${var.env}-bastion"

  vpc_zone_identifier = module.vpc.public_subnets

  desired_capacity          = "1"
  min_size                  = "1"
  max_size                  = "1"
  health_check_grace_period = "60"
  health_check_type         = "EC2"
  force_delete              = false
  wait_for_capacity_timeout = 0
  launch_configuration      = aws_launch_configuration.bastion.name

  enabled_metrics = [
    "GroupMinSize",
    "GroupMaxSize",
    "GroupDesiredCapacity",
    "GroupInServiceInstances",
    "GroupPendingInstances",
    "GroupStandbyInstances",
    "GroupTerminatingInstances",
    "GroupTotalInstances",
  ]

  # tags = [var.tags]

  lifecycle {
    create_before_destroy = true
  }
}