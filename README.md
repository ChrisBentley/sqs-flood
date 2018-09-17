# sqs-flood

Go utility for reading messages from a file and placing them on an SQS queue

Intended to be used with https://github.com/oscar-barlow/sqs-drain

## Installation

    go get github.com/ChrisBentley/sqs-flood

## Configuration

AWS credentials file ~/.aws/credentials with `aws_access_key_id` and `aws_secret_access_key`
AWS config file ~/.aws/config with required profiles

## Usage

Supply source file, destination URL endpoint and aws profile.

Region is defaulted to `eu-west-1` but can also be supplied.

If FIFO also include messageGroupId (optional)

e.g.

    ./sqs-flood -src "sqs-drain.json" -dest https://sqs.eu-west-1.amazonaws.com/<aws-account-number>/<sqs-queue-name> -profile sb-test

## License

The MIT License (MIT)

Copyright (c) 2016 Scott Barr

See [LICENSE.md](LICENSE.md)