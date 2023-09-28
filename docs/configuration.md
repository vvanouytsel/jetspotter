# Configuration

The following environment variables have to be set:

* SLACK_WEBHOOK_URL: The Webhook URL to send messages to Slack
* AIRCRAFT_TYPE: The [type of aircraft](./internal/aircraft/aircraft.go) you want to spot, if not specified all types will be spotted
* LOCATION_LATITUDE: The latitude coordinate (e.g.: `51.16951182347571`)
* LOCATION_LONGITUDE: The longitude coordinate of your location (e.g.: `5.470099273882526`)
* MAX_RANGE_KILOMETERS: The maximum range to spot aircraft from your location (default: `30`)
* MAX_AIRCRAFT_SLACK_MESSAGE: The amount of aircraft shown in a single slack message (default: `8`)
