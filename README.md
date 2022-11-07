[![Feature Flag, Remote Config and A/B Testing platform, Flagsmith](https://github.com/Flagsmith/flagsmith/raw/main/static-files/hero.png)](https://www.flagsmith.com/)

[Flagsmith](https://www.flagsmith.com/) is an open source, fully featured, Feature Flag and Remote Config service. Use
our hosted API, deploy to your own private cloud, or run on-premise.


# Flagsmith Terraform provider

This repository hosts the terraform provider that can be used to update Feature flags in Flagsmith.

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.17

## Building The Provider

1. Clone the repository
1. Enter the repository directory
1. Build the provider using the Go `install` command:

```shell
go install
```

## Adding Dependencies

This provider uses [Go modules](https://github.com/golang/go/wiki/Modules).
Please see the Go documentation for the most up to date information about using Go modules.

To add a new dependency `github.com/author/dependency` to your Terraform provider:

```shell
go get github.com/author/dependency
go mod tidy
```

Then commit the changes to `go.mod` and `go.sum`.

## Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (see [Requirements](#requirements) above).

To compile the provider, run `go install`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

To generate or update documentation, run `go generate`.

Running acceptance tests require following environment variables to be set.

* FLAGSMITH_MASTER_API_KEY
* FLAGSMITH_FEATURE_NAME
* FLAGSMITH_ENVIRONMENT_KEY
* FLAGSMITH_ENVIRONMENT_ID
* FLAGSMITH_FEATURE_ID
* FLAGSMITH_PROJECT_UUID

In order to run the full suite of Acceptance tests, run `make testacc`.

*Note:* Acceptance tests create real resources, and often cost money to run.

```shell
make testacc
```
