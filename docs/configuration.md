# Configuration

You can use environment variables to configure the application.
The supported parameters and their corresponding environment variables are listed below in the following format:

```go
// ENV_VARIABLE_NAME DEFAULT_VALUE
```

```go
{%
   include-markdown "snippets/config.snippet"
   comments=false
%}
```

You can also use a `config.yaml` file located in the same directory as the executable. Here is a sample:

```yaml
aircraftTypes: ALL
locationLatitude: 51.17348
locationLongitude: 5.45921
logNewPlanesToConsole: true
maxRangeKilometers: 20
metricsPort: 7070
slackWebhookUrl: https://hooks.slack.com/services/XXX/YYY/ZZZ
```
