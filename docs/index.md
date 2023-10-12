# Overview

JetSpotter is a simple program that queries the ADS-B API.
It is used to send notifications if a specified type of aircraft has been spotted within a specified range of a target location.
If one or more jets have been spotted, a notification is sent. The notification contains some metadata about the aircraft, a picture fetched from <a href="https://www.jetphotos.com" target="_blank">jetphotos.com</a> and a link to the <a href="https://globe.adsbexchange.com" target="_blank">ADS-B exchange page</a> of that aircraft.
A notification is only sent once for each aircraft. If the aircraft leaves your maximum configured range for at least 1 fetch iteration, a notification will be sent again as soon as it enters your maximum configured range.

## Notifications

Terminal output is always shown. Depending on the [configuration](configuration.md), notifications can also be sent via other media.

### Terminal

![Terminal output ](images/jetspotter-terminal-1.png)

### Slack

Slack notifications are sent if the `SLACK_WEBHOOK_URL` environment variable is defined.
Documentation how to set up notifications using incoming webhooks can be found in the [official slack documentation](https://api.slack.com/messaging/webhooks).

![Slack notfication ](images/jetspotter-slack-1.png)

### Discord

Discord notifications are sent if the `DISCORD_WEBHOOK_URL` environment variable is defined.
Documentation how to set up notifications using incoming webhooks can be found in the [official discord documentation](https://support.discord.com/hc/en-us/articles/228383668-Intro-to-Webhooks).

By default the color of the embed message is related to the altitude of the aircraft. The color scheme is the same as on the [ADS-B exchange map](https://globe.adsbexchange.com/). This feature can be disabled in the [configuration](configuration.md) to use the same static color for every embed message.

If the altitude color feature is enabled:
![Discord notfication ](images/jetspotter-discord-1.png)
![Altitude color scale ](images/jetspotter-color-scale.png)

If the altitude color feature is disabled:
![Discord notfication ](images/jetspotter-discord-2.png)

## Getting started

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

Send a slack notification if one or more aircraft are spotted

```bash
# Docker
docker run -e SLACK_WEBHOOK_URL=https://hooks.slack.com/services/XXX/YYY/ZZZ
 ghcr.io/vvanouytsel/jetspotter:latest

# Binary
SLACK_WEBHOOK_URL=https://hooks.slack.com/services/XXX/YYY/ZZZ ./jetspotter
```

Send a discord notification if one or more aircraft are spotted

```bash
# Docker
docker run -e DISCORD_WEBHOOK_URL=https://discord.com/api/webhooks/XXXXXX/YYYYYY
 ghcr.io/vvanouytsel/jetspotter:latest

# Binary
DISCORD_WEBHOOK_URL=https://discord.com/api/webhooks/XXXXXX/YYYYYY ./jetspotter
```

## Helm

Helm charts are available in the oci registry.
Configuration values can be found in the repository or via [artifact hub](https://artifacthub.io/packages/helm/jetspotter/jetspotter).

```bash
helm install -n jetspotter --create-namespace oci://ghcr.io/vvanouytsel/jetspotter-chart/jetspotter
```

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
