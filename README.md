# JetSpotter

JetSpotter is a simple program that queries the ADS-B API.
It is used to send notifications if a specified type of aircraft has been spotted within a specified range of a target location.

## Build

TODO

## Run

```bash
go run cmd/jetspotter/*
```

## Test

TODO

## Deploy

The following environment variables have to be set:

* SLACK_WEBHOOK_URL: The Webhook URL to send messages to Slack
