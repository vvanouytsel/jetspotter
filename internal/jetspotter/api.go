package jetspotter

import (
	"log"
	"net/http"
	"sync"

	"jetspotter/internal/auth"
	"jetspotter/internal/configuration"

	"github.com/gin-gonic/gin"
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

	// Set Gin to release mode in production
	gin.SetMode(gin.ReleaseMode)

	// Create a new Gin router
	router := gin.Default()

	// Create auth middleware
	basicAuth := auth.NewBasicAuth()

	// API routes
	router.GET("/api/aircraft", handleAircraftAPI)

	// Config API endpoint requires authentication
	router.GET("/api/config", basicAuth.Middleware(), handleConfigAPI)

	// Start HTTP server
	go func() {
		if err := router.Run(":" + listenPort); err != nil {
			log.Printf("API server error: %v", err)
		}
	}()
}

// handleAircraftAPI returns all currently spotted aircraft as JSON
func handleAircraftAPI(c *gin.Context) {
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
	config, err := configuration.GetConfig()
	if err != nil {
		log.Printf("Error getting config for API: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get configuration"})
		return
	}

	outputs, err := CreateAircraftOutput(filteredAircraft, config, true)
	if err != nil {
		log.Printf("Error creating aircraft output for API: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process aircraft data"})
		return
	}

	c.JSON(http.StatusOK, outputs)
}

// handleConfigAPI returns the application configuration as JSON
func handleConfigAPI(c *gin.Context) {
	// This endpoint is now protected by the auth middleware
	c.JSON(http.StatusOK, Config)
}
