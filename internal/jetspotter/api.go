package jetspotter

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"jetspotter/internal/configuration"
)

// SpottedAircraft keeps track of all currently spotted aircraft
var SpottedAircraft struct {
	sync.Mutex
	Aircraft []Aircraft
}

// Config holds the application configuration for API access
var Config configuration.Config

// SetupAPI sets up the API endpoints for the web server
func SetupAPI(listenPort string, config configuration.Config) {
	log.Printf("Serving API on port %s and path /api", listenPort)

	// Store the configuration for API access
	Config = config

	// Create API endpoints
	http.HandleFunc("/api/aircraft", handleAircraftAPI)
	http.HandleFunc("/api/config", handleConfigAPI)

	// Start HTTP server
	go func() {
		err := http.ListenAndServe(":"+listenPort, nil)
		if err != nil {
			log.Printf("API server error: %v", err)
		}
	}()
}

// handleAircraftAPI returns all currently spotted aircraft as JSON
func handleAircraftAPI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get the current spotted aircraft
	SpottedAircraft.Lock()
	airplanes := ConvertAircraftToOutput(SpottedAircraft.Aircraft)
	SpottedAircraft.Unlock()

	// Return as JSON
	json.NewEncoder(w).Encode(airplanes)
}

// handleConfigAPI returns the application configuration as JSON
func handleConfigAPI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Return the configuration as JSON
	json.NewEncoder(w).Encode(Config)
}
