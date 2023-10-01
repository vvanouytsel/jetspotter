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
![JetSpotter slack notfication ](images/jetspotter-slack-2.png)

## Releases

Releases can be found on the [GitHub repository](https://github.com/vvanouytsel/jetspotter/releases).

## Container images

[Container images](https://github.com/vvanouytsel/jetspotter/pkgs/container/jetspotter) are also created for each release.

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

## Examples

Run jetspotter without extra parameters.

```bash
# Docker
docker run ghcr.io/vvanouytsel/jetspotter:latest

# Binary
./jetspotter
```

Only show F16 and A400 aircraft within 100 kilometers of Kleine-Brogel airbase.

```bash
# Docker
docker run -e LOCATION_LATITUDE=51.1697898378895 -e LOCATION_LONGITUDE=5.470114381971933 -e AIRCRAFT_TYPES=F16,A400 -e MAX_RANGE_KILOMETERS=100 ghcr.io/vvanouytsel/jetspotter:latest


# Binary
LOCATION_LATITUDE=51.1697898378895 LOCATION_LONGITUDE=5.470114381971933 AIRCRAFT_TYPES=F16,A400 MAX_RANGE_KILOMETERS=100 ./jetspotter
```
