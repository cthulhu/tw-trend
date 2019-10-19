resource "aws_ssm_parameter" "consumer-key" {
  name        = "/${var.name}/${var.env}/CONSUMERKEY"
  type        = "SecureString"
  value       = "${var.twitter_consumer_key}"
  tags = var.tags
}

resource "aws_ssm_parameter" "consumer-secret" {
  name        = "/${var.name}/${var.env}/CONSUMERSECRET"
  type        = "SecureString"
  value       = "${var.twitter_consumer_secret}"
  tags = var.tags
}

resource "aws_ssm_parameter" "access-token" {
  name        = "/${var.name}/${var.env}/ACCESSTOKEN"
  type        = "SecureString"
  value       = "${var.twitter_access_token}"
  tags = var.tags
}


resource "aws_ssm_parameter" "access-secret" {
  name        = "/${var.name}/${var.env}/ACCESSSECRET"
  type        = "SecureString"
  value       = "${var.twitter_access_secret}"
  tags = var.tags
}
