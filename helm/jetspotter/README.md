# Overview

[![GitHub repository](https://img.shields.io/badge/GitHub-jetspotter-green?logo=github)](https://github.com/vvanouytsel/jetspotter)

Jetspotter is a simple program that queries the ADS-B API.
It is used to send notifications if a specified type of aircraft has been spotted within a specified range of a target location.
If one or more jets have been spotted, a notification is sent. The notification contains some metadata about the aircraft, a picture fetched from <a href="https://www.jetphotos.com" target="_blank">jetphotos.com</a> and a link to the <a href="https://globe.adsbexchange.com" target="_blank">ADS-B exchange page</a> of that aircraft.
A notification is only sent once for each aircraft. If the aircraft leaves your maximum configured range for at least 1 fetch iteration, a notification will be sent again as soon as it enters your maximum configured range.
