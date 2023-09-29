# Overview

JetSpotter is a simple program that queries the ADS-B API.
It is used to send notifications if a specified type of aircraft has been spotted within a specified range of a target location.
If one or more jets have been spotted, a slack notification is sent. The slack notification contains some metadata about the aircraft, a picture fetched from <a href="https://www.planespotting.be" target="_blank">planespotting.be</a> and a link to the <a href="https://globe.adsbexchange.com" target="_blank">ADS-B exchange page</a> of that aircraft.

## Notifications

Terminal output is always shown. Depending on the [configuration](configuration.md), notifications can also be sent via other media.

### Terminal

![JetSpotter CLI output ](images/jetspotter-cli.png)

### Slack

![JetSpotter slack notfication ](images/jetspotter-slack.png)

## Build

```bash
go build -o jetspotter -ldflags "-linkmode external -extldflags -static" cmd/jetspotter/jetspotter.go
```

## Run

```bash
go run cmd/jetspotter/*
```

## Test

```bash
make test
```

## Documentation

Documentation for the configuration parameters can be generated.

```bash
make doc
```
