# Twitter Trending

The application has to perform the following actions:

- Connect to a Twitter stream;

- Collect and store data of the latest 72h;

- Segment/Analyse the data in order to identify what is trending in the area of Amsterdam NL

- Expose the identified trends in a JSON feed.

## Getting Started



### Prerequisites

* Docker - to build and run locally
* github access to the repository in order to make changes
* terraform 0.12+
* Active access to AWS:
  * env variables AWS_ACCESS_KEY_ID/AWS_SECRET_ACCESS_KEY if user based access
  * env variables AWS_ACCESS_KEY_ID/AWS_SECRET_ACCESS_KEY/AWS_SESSION_TOKEN/AWS_SECURITY_TOKEN if role based access

## Infrastructure

To create infrastructure:

```
cd infra/aws && terraform apply
```

## Building

To build go to github https://github.com/cthulhu/tw-trend and press play on github actions. It will also run the tests

## Deployment

To deploy

pipeline???

```
cd infra/aws && terraform taint ... && terraform apply
```

## Running the tests

Docker has multiple stages one of them is to run tests as well. In order to run tests just build the container.


## Built With

* [Golang](https://golang.org/) - The language used
* [GoTwitter](https://github.com/dghubble/go-twitter/) - Twitter API sdk
* [Terraform](https://www.terraform.io/) - IaaC

* AWS infra:
   * AWS Fargate (https://aws.amazon.com/fargate/)
   * S3 (https://aws.amazon.com/s3/)
   * VPC and networking (https://aws.amazon.com/vpc/)
   * Alb (https://aws.amazon.com/elasticloadbalancing/)

## Author

* **Stanislav O. Pogrebnyak** - aka [Cthulhu](https://github.com/cthulhu)
