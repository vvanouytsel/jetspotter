# JetSpotter

JetSpotter is a simple program that queries the ADS-B API.
It is used to send notifications if a specified type of aircraft has been spotted within a specified range of a target location.
If one or more jets have been spotted, a slack notification is sent. The slack notification contains some metadata about the aircraft, a picture fetched from [planespotting.be](www.planespotting.be) and a link to the [ADS-B exchange page](https://globe.adsbexchange.com) of that aircraft.

![JetsSpotter slack notfication ](image/../images/jetspotter-slack.png)

## Build

Pushing to master triggers the release [workflow](.github/workflows/release.yaml).

## Run

```bash
go run cmd/jetspotter/*
```

## Test

```bash
go test -buildvcs=false ./internal/... -p 1 --count=1
```
