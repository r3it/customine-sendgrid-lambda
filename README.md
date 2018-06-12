Customine Sendmail Lambda
====

Simple send email Lambda function using SendGrid.

## Description

## Requirement

You must need SendGrid account. see https://sendgrid.com/

Requirements management depend on dep.

```
go get -u github.com/golang/dep/cmd/dep
```

**used libraries**

* https://github.com/sendgrid/sendgrid-go

## Build and Deploy with AWS SAM

```
$ GOOS=linux GOARCH=amd64 go build -o build/sendMailLambda

$ aws cloudformation package \
    --profile your_profile_name \
    --template-file template.yml \
    --s3-bucket your_bucket \
    --region your_region \
    --s3-prefix your_bucket_prefix \
    --output-template-file .template.yml
$ aws cloudformation deploy \
    --profile your_profile_name \
    --template-file .template.yml \
    --capabilities CAPABILITY_IAM \
    --stack-name yourStackName

$ aws cloudformation describe-stack-events \
    --profile your_profile_name \
    --stack-name yourStackName
```

## Setup

Please setup environment variables on your Lambda function.

| SENDGRID_API_KEY | your API Key |
| SENDER_NAME | sender name |
| SENDER_ADDRESS | sender email address |

## Usage

Request parameters:

```
{
    "subject": "mail subject",
    "to_name": "recipient name",
    "to_address": "recipient email",
    "body_text": "text mail body"
}
```

Response parameters:

```
{
    "status_code": "..."
}
```

## License

Apache License

## Author


[Koichiro Nishijima](https://github.com/k-nishijima/) / [R3 institute](https://www.r3it.com/)
