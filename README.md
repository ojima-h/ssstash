# ssstash

`ssstash` is a very simple tool that manages sensitive information in S3.
All entities are encrypted at client side using KMS before sending to S3 (see [Protecting Data Using Client-Side Encryption](http://docs.aws.amazon.com/AmazonS3/latest/dev/UsingClientSideEncryption.html)).

This is a S3 fork of [credstash](https://github.com/fugue/credstash). Please refer to it for basic concepts.

## Install

```
go get github.com/ojima-h/ssstash
```

or

Download from https://github.com/ojima-h/ssstash/releases

```
curl -L https://github.com/ojima-h/ssstash/releases/download/v0.0.1/ssstash-0.0.1.linux-amd64.gz | zcat > ssstash
chmod +x ssstash
```

## Usage

```
NAME:
   ssstash - A new cli application

USAGE:
   ssstash [global options] command [command options] [arguments...]

VERSION:
   0.0.1

COMMANDS:
     list, ls    List saved credentials
     put         Save the credential in S3
     get         Get the credential from S3
     delete, rm  Delete the entry
     help, h     Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help
   --version, -v  print the version
```

## Configurations

The following Environment Variables can be used to configure default options.

- `SSSTASH_S3_BUCKET` -- S3 bucket where credentials are saved.
- `SSSTASH_S3_PREFIX` -- S3 prefix under which credentials are saved.
- `SSSTASH_AWS_PROFILE` -- Profile name in your .aws credential file.
- `SSSTASH_KMS_KEY_ARN` -- KMS Key ARN to encrypt/decrypt credentials.

AWS CLI environment variables (`AWS_ACCESS_KEY_ID`, ... etc) are also available.
ref. [Configuring the AWS CLI » Environment Variables](http://docs.aws.amazon.com/cli/latest/userguide/cli-environment.html)

## Example

Create a new credential:

```
ssstash put passward very-very-secret --key arn:aws:kms:ap-east-1:000000000000:key/xxxx
```

Fetch the credential:

```
ssstash get password
```

List saved credentials:

```
ssstash ls
```

Delete the credential:

```
ssstash rm password
```

## credstash vs. ssstash

`credstash` uses DynamoDB to store credentials while `ssstash` uses S3.
Each of these has its own advantages.

The advantages to use DynamoDB is described [here](https://github.com/fugue/credstash#4-why-dynamodb-for-the-credential-store-why-not-s3).

The good points of S3 are:

* You don't have to care about the capacity.
* S3 is cheaper than DynamoDB in general.
* Many AWS SDKs supports S3 client-side encryption
  (see [AWS SDK Support for Amazon S3 Client-Side Encryption](http://docs.aws.amazon.com/general/latest/gr/aws_sdk_cryptography.html)).
  No special libraries are needed to fetch saved credentials.
* You can control access permissions for each entries using S3 Bucket Policy.
* S3 supports _Versioning_.