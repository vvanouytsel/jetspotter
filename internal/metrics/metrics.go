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
	[]string{"type", "description", "military"},
)

var aircraftAltitude = promauto.NewHistogramVec(prometheus.HistogramOpts{
	Name:    "jetspotter_aircraft_altitude_feet",
	Help:    "The altitude the jet is flying at.",
	Buckets: []float64{0, 2500, 5000, 10000, 15000, 20000, 25000, 30000, 35000, 40000},
},
	[]string{"type", "description", "military"},
)

// IncrementMetrics handles the metrics that need to be incremented
func IncrementMetrics(aircrafType, description, military string, altitude float64) {
	go func() {
		aircraftSpotted.WithLabelValues(aircrafType, description, military).Inc()
		aircraftAltitude.WithLabelValues(aircrafType, description, military).Observe(altitude)
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
