# Overview

[![GitHub repository](https://img.shields.io/badge/GitHub-jetspotter-green?logo=github)](https://github.com/vvanouytsel/jetspotter)

Jetspotter is a simple program that queries the ADS-B API.
It is used to send notifications if a specified type of aircraft has been spotted within a specified range of a target location.
If one or more jets have been spotted, a notification is sent. The notification contains some metadata about the aircraft, a picture fetched from <a href="https://www.planespotters.net" target="_blank">planespotters.net</a> and a link to <a href="https://globe.airplanes.live" target="_blank">track the aircraft</a>.
A notification is only sent once for each aircraft. If the aircraft leaves your maximum configured range for at least 1 fetch iteration, a notification will be sent again as soon as it enters your maximum configured range.
