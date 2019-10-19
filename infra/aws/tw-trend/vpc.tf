module "vpc" {
  source = "terraform-aws-modules/vpc/aws"
  version = "2.17.0"

  name = "${var.name}-${var.env}-vpc"

  cidr = "10.10.0.0/16"

  azs             = ["eu-west-1a", "eu-west-1b"]
  public_subnets  = ["10.10.1.0/24", "10.10.2.0/24"]
  private_subnets = ["10.10.11.0/24", "10.10.12.0/24"]

  enable_nat_gateway = true

  tags = var.tags
}
