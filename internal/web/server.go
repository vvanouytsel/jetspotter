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
	"jetspotter/internal/version"
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
	DataReadyChan <-chan bool // Channel to signal when data is ready
}

// Server represents the web frontend server
type Server struct {
	config        Config
	router        *http.ServeMux
	isDataReady   bool
	dataReadyMux  http.Handler
	pendingRoutes http.Handler
}

// NewServer creates a new web frontend server
func NewServer(config Config) *Server {
	pendingMux := http.NewServeMux()
	pendingMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(`<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Jetspotter - Starting</title>
    <meta http-equiv="refresh" content="5">
    <style>
        :root {
            --primary-color: #3498db;
            --secondary-color: #2c3e50;
            --accent-color: #e74c3c;
            --background-color: #f5f7fa;
            --card-color: #ffffff;
            --text-color: #333333;
        }
        
        @media (prefers-color-scheme: dark) {
            :root {
                --primary-color: #e0e0e0;
                --secondary-color: #cccccc;
                --accent-color: #ff6b6b;
                --background-color: #121212;
                --card-color: #1e1e1e;
                --text-color: #e0e0e0;
            }
        }
        
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }
        
        body {
            font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
            background-color: var(--background-color);
            color: var(--text-color);
            line-height: 1.6;
            display: flex;
            flex-direction: column;
            min-height: 100vh;
        }
        
        .container {
            max-width: 1200px;
            margin: 0 auto;
            padding: 20px;
            flex: 1;
            display: flex;
            flex-direction: column;
            justify-content: center;
            align-items: center;
        }
        
        header {
            text-align: center;
            margin-bottom: 2rem;
            padding-bottom: 1.5rem;
            width: 100%;
            border-bottom: 1px solid rgba(125, 125, 125, 0.2);
        }
        
        header h1 {
            color: var(--secondary-color);
            font-size: 2.5rem;
            margin-bottom: 0.5rem;
        }
        
        .loading-container {
            display: flex;
            flex-direction: column;
            align-items: center;
            justify-content: center;
            background-color: var(--card-color);
            border-radius: 12px;
            padding: 3rem;
            box-shadow: 0 8px 24px rgba(0, 0, 0, 0.1);
            width: 100%;
            max-width: 600px;
            margin: 0 auto;
            position: relative;
            overflow: hidden;
        }
        
        .scanning-animation {
            position: relative;
            width: 160px;
            height: 160px;
            margin-bottom: 30px;
        }
        
        .radar-circle {
            position: absolute;
            top: 0;
            left: 0;
            width: 100%;
            height: 100%;
            border-radius: 50%;
            border: 3px solid var(--secondary-color);
            box-shadow: 0 0 20px rgba(0, 0, 0, 0.2);
            animation: pulse 2s infinite ease-in-out;
        }
        
        @keyframes pulse {
            0% { transform: scale(0.97); opacity: 0.8; }
            50% { transform: scale(1.03); opacity: 1; }
            100% { transform: scale(0.97); opacity: 0.8; }
        }
        
        .radar-circle::before, .radar-circle::after {
            content: '';
            position: absolute;
            top: 50%;
            left: 50%;
            transform: translate(-50%, -50%);
            border-radius: 50%;
            border: 2px solid var(--secondary-color);
            opacity: 0.7;
        }
        
        .radar-circle::before {
            width: 70%;
            height: 70%;
            animation: ripple 3s infinite ease-out;
        }
        
        .radar-circle::after {
            width: 40%;
            height: 40%;
            animation: ripple 3s infinite ease-out 1s;
        }
        
        @keyframes ripple {
            0% { opacity: 0.8; transform: translate(-50%, -50%) scale(0.9); }
            50% { opacity: 0.4; transform: translate(-50%, -50%) scale(1.1); }
            100% { opacity: 0.8; transform: translate(-50%, -50%) scale(0.9); }
        }
        
        .radar-sweep {
            position: absolute;
            top: 50%;
            left: 50%;
            width: 52%;
            height: 3px;
            background-color: var(--secondary-color);
            transform-origin: left center;
            animation: radar-sweep 4s linear infinite;
            box-shadow: 0 0 10px rgba(0, 0, 0, 0.3);
            border-radius: 3px;
        }
        
        .radar-sweep::after {
            content: '';
            position: absolute;
            top: 0;
            right: 0;
            width: 40px;
            height: 100%;
            background: linear-gradient(to right, var(--secondary-color), transparent);
            border-radius: 3px;
        }
        
        @keyframes radar-sweep {
            0% { transform: rotate(0deg); }
            100% { transform: rotate(360deg); }
        }
        
        .aircraft-dots {
            position: absolute;
            width: 8px;
            height: 8px;
            border-radius: 50%;
            background-color: var(--accent-color);
            opacity: 0;
            animation: dot-appear 8s infinite linear;
        }
        
        .aircraft-dot-1 {
            top: 35%;
            left: 70%;
            animation-delay: 1s;
        }
        
        .aircraft-dot-2 {
            top: 60%;
            left: 30%;
            animation-delay: 3.5s;
        }
        
        .aircraft-dot-3 {
            top: 20%;
            left: 40%;
            animation-delay: 6s;
        }
        
        @keyframes dot-appear {
            0% { opacity: 0; transform: scale(0); }
            5% { opacity: 1; transform: scale(1); }
            15% { opacity: 1; transform: scale(1); }
            20% { opacity: 0; transform: scale(0); }
            100% { opacity: 0; transform: scale(0); }
        }
        
        .loading-text {
            text-align: center;
            max-width: 80%;
        }
        
        .loading-title {
            font-size: 1.8rem;
            font-weight: 600;
            margin-bottom: 1rem;
            color: var(--secondary-color);
        }
        
        .loading-message {
            font-size: 1.1rem;
            margin-bottom: 1.5rem;
            color: var(--text-color);
            opacity: 0.9;
        }
        
        .loading-progress {
            width: 100%;
            height: 4px;
            background-color: rgba(125, 125, 125, 0.2);
            border-radius: 2px;
            overflow: hidden;
            position: relative;
            margin-top: 1.5rem;
        }
        
        .loading-progress-bar {
            position: absolute;
            height: 100%;
            background-color: var(--secondary-color);
            width: 30%;
            border-radius: 2px;
            animation: progress-animation 3s infinite ease-in-out;
        }
        
        @keyframes progress-animation {
            0% { width: 0%; left: 0; }
            50% { width: 30%; }
            100% { left: 100%; width: 0%; }
        }
        
        .refresh-note {
            font-size: 0.9rem;
            color: var(--text-color);
            opacity: 0.7;
            margin-top: 1.5rem;
            text-align: center;
        }
        
        footer {
            text-align: center;
            margin-top: 2rem;
            padding: 1rem 0;
            opacity: 0.7;
            font-size: 0.9rem;
        }
        
        @media (prefers-color-scheme: dark) {
            .radar-circle {
                border-color: var(--accent-color);
                box-shadow: 0 0 30px rgba(231, 76, 60, 0.2);
            }
            
            .radar-circle::before, .radar-circle::after {
                border-color: var(--accent-color);
            }
            
            .radar-sweep {
                background-color: var(--accent-color);
                box-shadow: 0 0 15px rgba(231, 76, 60, 0.4);
            }
            
            .radar-sweep::after {
                background: linear-gradient(to right, var(--accent-color), transparent);
            }
            
            .loading-title {
                color: var(--accent-color);
            }
            
            .loading-progress-bar {
                background-color: var(--accent-color);
            }
        }
    </style>
</head>
<body>
    <div class="container">
        <header>
            <h1>Jetspotter</h1>
        </header>
        
        <div class="loading-container">
            <div class="scanning-animation">
                <div class="radar-circle"></div>
                <div class="radar-sweep"></div>
                <div class="aircraft-dots aircraft-dot-1"></div>
                <div class="aircraft-dots aircraft-dot-2"></div>
                <div class="aircraft-dots aircraft-dot-3"></div>
            </div>
            
            <div class="loading-text">
                <div class="loading-title">Jetspotter is starting...</div>
                <div class="loading-message">Initializing systems and waiting for aircraft data to become available.</div>
                <div class="loading-progress">
                    <div class="loading-progress-bar"></div>
                </div>
                <div class="refresh-note">This page will refresh automatically every 5 seconds.</div>
            </div>
        </div>
    </div>
    
    <footer>
        <div>&copy; Jetspotter</div>
    </footer>
</body>
</html>`))
	})

	mainMux := http.NewServeMux()

	server := &Server{
		config:        config,
		router:        mainMux,
		isDataReady:   false,
		dataReadyMux:  pendingMux,
		pendingRoutes: pendingMux,
	}

	server.setupRoutes()

	// Set up a goroutine to wait for the data ready signal
	if config.DataReadyChan != nil {
		go server.waitForData()
	}

	return server
}

// waitForData waits for the data ready signal
func (s *Server) waitForData() {
	log.Println("Web UI: Waiting for aircraft data to be available...")
	<-s.config.DataReadyChan
	s.isDataReady = true
	s.dataReadyMux = s.router // Switch to the main router once data is ready
	log.Println("Web UI: Aircraft data is now available, serving full interface")
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

	// Configuration page
	s.router.HandleFunc("/config", s.handleConfig)

	// API proxy - forwards to the actual API for AJAX requests
	s.router.HandleFunc("/api/aircraft", s.handleAPIProxy)

	// API endpoint for configuration
	s.router.HandleFunc("/api/config", s.handleAPIConfigProxy)

	// Version information endpoint
	s.router.HandleFunc("/api/version", s.handleVersion)
}

// ServeHTTP makes Server implement the http.Handler interface
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.dataReadyMux.ServeHTTP(w, r)
}

// Start runs the web server
func (s *Server) Start() error {
	addr := fmt.Sprintf(":%d", s.config.ListenPort)
	log.Printf("Starting web frontend at http://localhost%s", addr)
	return http.ListenAndServe(addr, s)
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

// handleConfig serves the configuration page
func (s *Server) handleConfig(w http.ResponseWriter, r *http.Request) {
	tmplFS, err := fs.Sub(content, "templates")
	if err != nil {
		http.Error(w, "Failed to load templates", http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFS(tmplFS, "config.html")
	if err != nil {
		http.Error(w, "Failed to parse template", http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"Title": "Jetspotter Configuration",
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

// handleAPIConfigProxy proxies requests to the backend API for configuration
func (s *Server) handleAPIConfigProxy(w http.ResponseWriter, r *http.Request) {
	// Forward the request to the actual API
	resp, err := http.Get(s.config.APIEndpoint + "/api/config")
	if err != nil {
		http.Error(w, "Failed to fetch config data from API", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Set the content type
	w.Header().Set("Content-Type", "application/json")

	// Copy the response body to the output
	var config interface{}
	if err := json.NewDecoder(resp.Body).Decode(&config); err != nil {
		http.Error(w, "Failed to parse API config response", http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(config); err != nil {
		http.Error(w, "Failed to encode config response", http.StatusInternalServerError)
	}
}

// handleVersion serves the application version information
func (s *Server) handleVersion(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(version.GetFullVersionInfo()); err != nil {
		http.Error(w, "Failed to encode version information", http.StatusInternalServerError)
	}
}
