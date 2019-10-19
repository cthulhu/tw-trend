resource "aws_iam_role" "ecs-service" {
  name = "${var.name}-${var.env}-ecs-role"

  assume_role_policy = <<EOF
{
  "Version": "2008-10-17",
  "Statement": [
	{
	  "Sid": "",
	  "Effect": "Allow",
	  "Principal": {
		"Service": "ecs.amazonaws.com"
	  },
	  "Action": "sts:AssumeRole"
	}
  ]
}
EOF
}

resource "aws_iam_role_policy_attachment" "ecs-service" {
  role       = "${aws_iam_role.ecs-service.name}"
  policy_arn = "arn:aws:iam::aws:policy/service-role/AmazonEC2ContainerServiceRole"
}

resource "aws_ecs_service" "ecs-service" {
  name            = "${var.name}-${var.env}-service"
  cluster         = aws_ecs_cluster.ecs_cluster.name
  task_definition = aws_ecs_task_definition.ecs_task_definition.arn
  desired_count   = var.ecs_desired_count
  iam_role        = aws_iam_role.ecs-service.arn

  deployment_maximum_percent         = var.deployment_maximum_percent
  deployment_minimum_healthy_percent = var.deployment_minimum_healthy_percent

  load_balancer {
    target_group_arn = module.alb.target_group_arns[0]
    container_name   = "tw-trend-container"
    container_port   = var.container_port
  }
}
