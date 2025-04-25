package web

import (
	"embed"
	"encoding/json"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"time"

	"jetspotter/internal/jetspotter"
)

// Templates and static assets embedded in the binary
//
//go:embed templates/* static/*
var content embed.FS

// Config holds the web server configuration
type Config struct {
	ListenPort    int
	APIEndpoint   string
	RefreshPeriod time.Duration
}

// Server represents the web frontend server
type Server struct {
	config Config
	router *http.ServeMux
}

// NewServer creates a new web frontend server
func NewServer(config Config) *Server {
	server := &Server{
		config: config,
		router: http.NewServeMux(),
	}

	server.setupRoutes()
	return server
}

// setupRoutes configures the HTTP routes
func (s *Server) setupRoutes() {
	// Serve static files
	staticFS, err := fs.Sub(content, "static")
	if err != nil {
		log.Fatalf("Failed to create static file server: %v", err)
	}
	s.router.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(staticFS))))

	// Homepage
	s.router.HandleFunc("/", s.handleIndex)

	// API proxy - forwards to the actual API for AJAX requests
	s.router.HandleFunc("/api/aircraft", s.handleAPIProxy)
}

// Start runs the web server
func (s *Server) Start() error {
	addr := fmt.Sprintf(":%d", s.config.ListenPort)
	log.Printf("Starting web frontend at http://localhost%s", addr)
	return http.ListenAndServe(addr, s.router)
}

// handleIndex serves the main page
func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	tmplFS, err := fs.Sub(content, "templates")
	if err != nil {
		http.Error(w, "Failed to load templates", http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFS(tmplFS, "index.html")
	if err != nil {
		http.Error(w, "Failed to parse template", http.StatusInternalServerError)
		return
	}

	// Get scan radius from the API endpoint
	resp, err := http.Get(s.config.APIEndpoint + "/api/config")
	var scanRadius int = 30 // Default value

	if err == nil {
		defer resp.Body.Close()
		var configData map[string]interface{}
		if decodeErr := json.NewDecoder(resp.Body).Decode(&configData); decodeErr == nil {
			// Prefer MaxScanRangeKilometers if available, otherwise fall back to MaxRangeKilometers
			if scanRange, ok := configData["MaxScanRangeKilometers"].(float64); ok {
				scanRadius = int(scanRange)
			} else if radius, ok := configData["MaxRangeKilometers"].(float64); ok {
				scanRadius = int(radius)
			}
		}
	}

	data := map[string]interface{}{
		"Title":         "Jetspotter",
		"RefreshPeriod": int(s.config.RefreshPeriod.Seconds()),
		"ScanRadius":    scanRadius,
	}

	w.Header().Set("Content-Type", "text/html")
	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
	}
}

// handleAPIProxy proxies requests to the backend API
func (s *Server) handleAPIProxy(w http.ResponseWriter, r *http.Request) {
	// Forward the request to the actual API
	resp, err := http.Get(s.config.APIEndpoint + "/api/aircraft")
	if err != nil {
		http.Error(w, "Failed to fetch data from API", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Parse the response
	var aircraft []jetspotter.AircraftOutput
	if err := json.NewDecoder(resp.Body).Decode(&aircraft); err != nil {
		http.Error(w, "Failed to parse API response", http.StatusInternalServerError)
		return
	}

	// Return the response as JSON
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(aircraft); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
