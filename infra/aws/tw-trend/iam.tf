resource "aws_iam_server_certificate" "server-cert" {
  name_prefix      = "${var.name}-${var.env}-cert"
  certificate_body = file("my-aws-public.crt")
  private_key      = file("my-aws-private.key")

  lifecycle {
    create_before_destroy = true
  }
}