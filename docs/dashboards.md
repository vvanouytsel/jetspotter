# Dashboards

Metrics are exposed via `$ADDRESS:7070/metrics`. All `jetspotter` specific metrics are prefixed with `jetspotter_`. These are usefull if you want to monitor how much aircraft per type are spotted at your location.

```bash
❯ curl -s localhost:7070/metrics  | grep jetspotter_

# HELP jetspotter_aircraft_spotted_total The total number of spotted aircraft.
# TYPE jetspotter_aircraft_spotted_total counter
jetspotter_aircraft_spotted_total{type="AERMACCHI SF-260"} 2
jetspotter_aircraft_spotted_total{type="AIRBUS A-319"} 3
jetspotter_aircraft_spotted_total{type="AIRBUS A-320"} 3
jetspotter_aircraft_spotted_total{type="AIRBUS A-321"} 2
jetspotter_aircraft_spotted_total{type="AIRBUS A-380-800"} 1
jetspotter_aircraft_spotted_total{type="BOEING 737-700"} 1
jetspotter_aircraft_spotted_total{type="BOEING 737-800"} 2
jetspotter_aircraft_spotted_total{type="BOEING 747-400"} 1
jetspotter_aircraft_spotted_total{type="BOEING 767-300"} 1
jetspotter_aircraft_spotted_total{type="BOEING 777-200"} 2
jetspotter_aircraft_spotted_total{type="DASSAULT Falcon 7X"} 1
```

Grafana can be used to create fancy dashboards.

Show military aircraft.
![Dashboard](images/jetspotter-grafana-1.png)

Show all aircraft.
![Dashboard](images/jetspotter-grafana-2.png)

Show specific types of aircraft.
![Dashboard](images/jetspotter-grafana-3.png)
