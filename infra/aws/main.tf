
module "tw-rend-staging" {
  source = "./tw-trend"
  name   = "tw-trend"
  env    = "staging"
  tags = {
    Environment = "staging"
    Owner       = "spogrebnyak@xebia.com"
    Name        = "tw-trends"
  }
}
