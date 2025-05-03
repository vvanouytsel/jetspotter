# Jetspotter

![Release workflow](https://github.com/vvanouytsel/jetspotter/actions/workflows/release.yaml/badge.svg)
[![Latest release](https://img.shields.io/github/v/release/vvanouytsel/jetspotter)](https://github.com/vvanouytsel/jetspotter/releases)
[![GitHub Release Date - Published_At](https://img.shields.io/github/release-date/vvanouytsel/jetspotter)](https://github.com/vvanouytsel/jetspotter/releases)
[![GitHub deployments](https://img.shields.io/github/deployments/vvanouytsel/jetspotter/github-pages?label=Documentation&link=https%3A%2F%2Fvvanouytsel.github.io%2Fjetspotter%2F)](https://vvanouytsel.github.io/jetspotter/)
[![Artifact Hub](https://img.shields.io/endpoint?url=https://artifacthub.io/badge/repository/jetspotter)](https://artifacthub.io/packages/search?repo=jetspotter)  
[!["Buy Me A Coffee"](https://www.buymeacoffee.com/assets/img/custom_images/orange_img.png)](https://www.buymeacoffee.com/vvanouytsel)

Jetspotter is a simple program that queries the ADS-B API. It is used to send notifications if a specified type of aircraft has been spotted within a specified range of a target location. If one or more jets have been spotted, a notification is sent. The notification contains some metadata about the aircraft, a picture fetched from planespotters.net and a link to track the aircraft. A notification is only sent once for each aircraft. If the aircraft leaves your maximum configured range for at least 1 fetch iteration, a notification will be sent again as soon as it enters your maximum configured range.

## [Documentation](https://vvanouytsel.github.io/jetspotter/)

Please have a look at the [documentation](https://vvanouytsel.github.io/jetspotter/) for installation and configuration steps.

## Demo

There is a [demo](https://bru.jetspotter.vvanouytsel.dev/) of the web interface available that shows aircraft in the vicinity of the [Brussels airport](https://bru.jetspotter.vvanouytsel.dev/).

## Screenshots

### Web UI
![Jetspotter UI](docs/images/jetspotter-ui-1.png)

### Discord Notifications
![Discord Notifications](docs/images/jetspotter-discord-1.png)

### Slack Notifications
![Slack Notifications](docs/images/jetspotter-slack-1.png)

### Ntfy Notifications
![Ntfy Notifications](docs/images/jetspotter-ntfy-1.png)

### Gotify Notifications
![Gotify Notifications](docs/images/jetspotter-gotify-1.png)

### Terminal Output
![Terminal Output](docs/images/jetspotter-terminal-1.png)

## Stargazers

[![Stargazers](https://starchart.cc/vvanouytsel/jetspotter.svg)](https://starchart.cc/vvanouytsel/jetspotter)
