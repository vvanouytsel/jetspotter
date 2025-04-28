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
	SpottedAircraft.Lock()
	defer SpottedAircraft.Unlock()

	// Filter out aircraft without registration before converting to output
	var filteredAircraft []Aircraft
	for _, ac := range SpottedAircraft.Aircraft {
		if ac.Registration != "" {
			filteredAircraft = append(filteredAircraft, ac)
		}
	}

	// Convert filtered aircraft to output format
	outputs := ConvertAircraftToOutput(filteredAircraft)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(outputs)
}

// handleConfigAPI returns the application configuration as JSON
func handleConfigAPI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Return the configuration as JSON
	json.NewEncoder(w).Encode(Config)
}
