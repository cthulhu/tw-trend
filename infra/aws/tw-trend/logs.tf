resource "aws_cloudwatch_log_group" "ecs" {
  name = "tf-ecs-group/ecs-agent-${var.name}-${var.env}"
}

resource "aws_cloudwatch_log_group" "app" {
  name = "tf-ecs-group/app-${var.name}-${var.env}"
}

resource "aws_cloudwatch_log_group" "instance" {
  name = "tf-ecs-group/instance-${var.name}-${var.env}"
}
