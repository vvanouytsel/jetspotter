# JetSpotter

JetSpotter is a simple program that queries the ADS-B API.
It is used to send notifications if a specified type of aircraft has been spotted within a specified range of a target location.

## Build

Pushing to master triggers the release [workflow](./workflows/release.yaml).

## Run

```bash
go run cmd/jetspotter/*
```

## Test

```bash
go test -buildvcs=false ./internal/... -p 1 --count=1
```

## Deploy

The following environment variables have to be set:

* SLACK_WEBHOOK_URL: The Webhook URL to send messages to Slack
