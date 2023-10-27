package metrics

import (
	"fmt"
	"jetspotter/internal/configuration"
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var aircraftSpotted = promauto.NewCounterVec(prometheus.CounterOpts{
	Name: "jetspotter_aircraft_spotted_total",
	Help: "The total number of spotted aircraft.",
},
	[]string{"type", "description"},
)

// IncrementAircraftSpotted increments the counter of spotted aircraft
func IncrementAircraftSpotted(aircrafType, description string) {
	go func() {
		aircraftSpotted.WithLabelValues(aircrafType, description).Add(1)
	}()
}

func HandleMetrics(config configuration.Config) error {
	path := "/metrics"
	port := config.MetricsPort
	http.Handle(path, promhttp.Handler())
	log.Printf("Serving metrics on port %s and path %s", port, path)
	err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
	if err != nil {
		return err
	}
	return nil
}
