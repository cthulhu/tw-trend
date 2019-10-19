variable "twitter_consumer_key" {}
variable "twitter_consumer_secret" {}
variable "twitter_access_token" {}
variable "twitter_access_secret" {}

module "tw-trend-staging" {
  source = "./tw-trend"
  name   = "tw-trend"
  env    = "staging"

  twitter_consumer_key = var.twitter_consumer_key
  twitter_consumer_secret = var.twitter_consumer_secret
  twitter_access_token = var.twitter_access_token
  twitter_access_secret = var.twitter_access_secret

  tags = {
    Environment = "staging"
    Owner       = "spogrebnyak@xebia.com"
    Name        = "tw-trends"
  }
}
