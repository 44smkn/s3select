# s3select

![tests](https://github.com/44smkn/s3select/actions/workflows/tests.yaml/badge.svg)
[![Go Reference](https://pkg.go.dev/badge/github.com/44smkn/s3select.svg)](https://pkg.go.dev/github.com/44smkn/s3select)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

`s3select` is a tool that simplifies the use of s3 select, a feature that allows you to retrieve a subset of S3 objects.

## Demo

This demo sets s3select configuration up and retrive data from s3 using s3select query.

![schreenshot of s3select query](https://raw.githubusercontent.com/44smkn/s3select/main/.github/images/s3select_screenshot.gif)

## Installation

To download the latest release, run:

```bash
curl -sL "https://github.com/44smkn/s3select/releases/latest/download/s3select_$(uname -s)_amd64.tar.gz" | tar xz -C /tmp
sudo mv /tmp/s3select /usr/local/bin
```

You will need to have AWS API credentials configured. What works for AWS CLI, should be sufficient. You can use [~/.aws/credentials file](https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-files.html) or [environment variables](https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-envvars.html#envvars-set).

## Basic Usage

To retrive subset of s3object, run:

```sh
s3select configure --profile alb-accesslog
s3select query -b $BUCKET -k $KEY_PREFIX -p alb-accesslog -e "SELECT s._9 as elb_status_code, s._13 as request FROM s3object s WHERE s._1 = 'https'"
```

The result of running the query will be output stdout.
