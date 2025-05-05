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

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"

	"jetspotter/internal/auth"
	"jetspotter/internal/configuration"
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
	SecureCookies bool        // Enable secure cookies (for HTTPS)
}

// Server represents the web frontend server
type Server struct {
	config           Config
	engine           *gin.Engine
	isDataReady      bool
	auth             *auth.BasicAuth
	jetspotterConfig *configuration.Config // Add jetspotter configuration
}

// NewServer creates a new web frontend server
func NewServer(config Config) *Server {
	// Set Gin to release mode in production
	gin.SetMode(gin.ReleaseMode)

	// Create a new gin engine
	engine := gin.New()

	// Use the recovery middleware
	engine.Use(gin.Recovery())

	// Generate a secure key for session encryption
	key := generateSecureKey()

	// Configure session middleware with a memory store (server-side sessions)
	// This is more secure than cookie store as session data is kept on the server
	store := memstore.NewStore(key)
	store.Options(sessions.Options{
		Path:     "/",
		MaxAge:   86400, // 1 day
		HttpOnly: true,
		Secure:   config.SecureCookies,
		SameSite: http.SameSiteLaxMode,
	})
	engine.Use(sessions.Sessions("jetspotter_session", store))

	// Get the jetspotter configuration
	jetspotterConfig, err := configuration.GetConfig()
	if err != nil {
		log.Printf("Warning: Failed to load jetspotter configuration: %v", err)
	}

	server := &Server{
		config:           config,
		engine:           engine,
		isDataReady:      false,
		auth:             auth.NewBasicAuth(),
		jetspotterConfig: &jetspotterConfig,
	}

	// Add custom logger middleware that skips static file requests
	engine.Use(func(c *gin.Context) {
		// Skip logging for static file requests
		if len(c.Request.URL.Path) > 8 && c.Request.URL.Path[:8] == "/static/" {
			c.Next()
			return
		}
		// Log other requests
		start := time.Now()
		path := c.Request.URL.Path
		c.Next()
		latency := time.Since(start)
		status := c.Writer.Status()
		log.Printf("[GIN] %d | %s | %s", status, latency, path)
	})

	// Configure routes
	server.setupRoutes()

	// Set up a handler for when data is not yet ready
	if config.DataReadyChan != nil {
		go server.waitForData(config.DataReadyChan)
	}

	return server
}

// generateSecureKey creates a random key for the cookie store
func generateSecureKey() []byte {
	// Use a fixed key for now to ensure sessions remain valid across restarts
	// In production, this should be stored securely and loaded from environment/config
	return []byte("jetspotter_secure_cookie_key_change_in_production")
}

// waitForData waits for the data ready signal
func (s *Server) waitForData(readyChan <-chan bool) {
	log.Println("Web UI: Waiting for aircraft data to be available...")
	<-readyChan
	s.isDataReady = true
	log.Println("Web UI: Aircraft data is now available, serving full interface")
}

// setupRoutes configures the HTTP routes
func (s *Server) setupRoutes() {
	// Create a filesystem for static files
	staticFS, err := fs.Sub(content, "static")
	if err != nil {
		log.Fatalf("Failed to create static file server: %v", err)
	}
	s.engine.StaticFS("/static", http.FS(staticFS))

	// Add specific route for favicon.ico
	s.engine.GET("/favicon.ico", func(c *gin.Context) {
		c.FileFromFS("images/favicon.ico", http.FS(staticFS))
	})

	// Authentication routes
	s.engine.GET("/login", s.handleLoginGet)
	s.engine.POST("/login", s.handleLoginPost)
	s.engine.GET("/logout", s.handleLogout)

	// Public routes
	s.engine.GET("/", s.handleIndex)
	s.engine.GET("/api/aircraft", s.handleAPIProxy)
	s.engine.GET("/api/version", s.handleVersion)

	// Protected routes using auth middleware
	protected := s.engine.Group("/")
	protected.Use(s.authRequired())
	{
		protected.GET("/config", s.handleConfig)
		protected.GET("/api/config", s.handleAPIConfigProxy)
	}

	// Add a NoRoute handler for 404 errors
	s.engine.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"message": "Page not found"})
	})
}

// authRequired is a middleware that checks if the user is authenticated
func (s *Server) authRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		user := session.Get("user")
		if user == nil {
			// User not logged in, redirect to login page
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
			return
		}
		// User is logged in, continue
		c.Next()
	}
}

// handleLoginGet handles GET requests to the login page
func (s *Server) handleLoginGet(c *gin.Context) {
	// Check if already logged in
	session := sessions.Default(c)
	if session.Get("user") != nil {
		c.Redirect(http.StatusFound, "/")
		return
	}

	// Get coordinates for the map link
	var latitude, longitude float64
	if s.jetspotterConfig != nil {
		// Location is a geodist.Coord struct, not a pointer
		latitude = s.jetspotterConfig.Location.Lat
		longitude = s.jetspotterConfig.Location.Lon
	}

	// Show login page
	showError := c.Query("error") == "true"
	c.HTML(http.StatusOK, "login.html", gin.H{
		"Title":     "Login - Jetspotter",
		"ShowError": showError,
		"Latitude":  latitude,
		"Longitude": longitude,
	})
}

// handleLoginPost handles POST requests to the login endpoint
func (s *Server) handleLoginPost(c *gin.Context) {
	// Get form data
	username := c.PostForm("username")
	password := c.PostForm("password")

	// Check credentials
	if username == s.auth.Username && password == s.auth.Password {
		// Valid credentials, create session
		session := sessions.Default(c)
		session.Set("user", username)
		session.Save()

		// Redirect to home page
		c.Redirect(http.StatusFound, "/")
		return
	}

	// Invalid credentials
	c.Redirect(http.StatusFound, "/login?error=true")
}

// handleLogout logs the user out
func (s *Server) handleLogout(c *gin.Context) {
	// Clear the session
	session := sessions.Default(c)
	session.Clear()
	session.Save()

	// Redirect to login page
	c.Redirect(http.StatusFound, "/login")
}

// handleIndex serves the main page
func (s *Server) handleIndex(c *gin.Context) {
	// Check if data is ready
	if !s.isDataReady && s.config.DataReadyChan != nil {
		s.serveLoadingPage(c)
		return
	}

	// Get scan radius directly from the configuration
	var scanRadius int = 30 // Default value
	var latitude, longitude float64
	if s.jetspotterConfig != nil {
		// Prefer MaxScanRangeKilometers if available
		scanRadius = s.jetspotterConfig.MaxScanRangeKilometers
		// Get location coordinates for the map link
		// Location is a geodist.Coord struct, not a pointer, so we can access its fields directly
		latitude = s.jetspotterConfig.Location.Lat
		longitude = s.jetspotterConfig.Location.Lon
	}

	// Check if user is logged in - safely access session
	isLoggedIn := false
	var username string

	// Safely access session data with error handling
	session := sessions.Default(c)
	user := session.Get("user")
	if user != nil {
		isLoggedIn = true
		var ok bool
		username, ok = user.(string)
		if !ok {
			// If user value is not a string, handle the error
			log.Printf("Warning: user session value is not a string type")
			// Clear invalid session
			session.Clear()
			session.Save()
			isLoggedIn = false
		}
	}

	// Render template
	c.HTML(http.StatusOK, "index.html", gin.H{
		"Title":         "Jetspotter",
		"RefreshPeriod": int(s.config.RefreshPeriod.Seconds()),
		"ScanRadius":    scanRadius,
		"IsLoggedIn":    isLoggedIn,
		"Username":      username,
		"Latitude":      latitude,
		"Longitude":     longitude,
	})
}

// serveLoadingPage serves the loading page while waiting for data
func (s *Server) serveLoadingPage(c *gin.Context) {
	c.HTML(http.StatusOK, "loading.html", gin.H{
		"Title": "Jetspotter - Starting",
	})
}

// handleConfig serves the configuration page
func (s *Server) handleConfig(c *gin.Context) {
	// User is already authenticated by middleware
	session := sessions.Default(c)
	user := session.Get("user")

	// Safe access to username
	var username string
	if user != nil {
		var ok bool
		username, ok = user.(string)
		if !ok {
			log.Printf("Warning: user session value is not a string type")
			// Redirect to login page if session data is invalid
			c.Redirect(http.StatusFound, "/login")
			return
		}
	} else {
		// This shouldn't happen due to authRequired middleware, but handle it anyway
		c.Redirect(http.StatusFound, "/login")
		return
	}

	// Get location coordinates for the map link
	var latitude, longitude float64
	if s.jetspotterConfig != nil {
		// Location is a geodist.Coord struct, not a pointer
		latitude = s.jetspotterConfig.Location.Lat
		longitude = s.jetspotterConfig.Location.Lon
	}

	c.HTML(http.StatusOK, "config.html", gin.H{
		"Title":      "Jetspotter Configuration",
		"IsLoggedIn": true,
		"Username":   username,
		"Latitude":   latitude,
		"Longitude":  longitude,
	})
}

// handleAPIProxy proxies requests to the backend API
func (s *Server) handleAPIProxy(c *gin.Context) {
	// Forward the request to the actual API
	resp, err := http.Get(s.config.APIEndpoint + "/api/aircraft")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch data from API"})
		return
	}
	defer resp.Body.Close()

	// Parse the response
	var aircraft []jetspotter.Aircraft
	if err := json.NewDecoder(resp.Body).Decode(&aircraft); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse API response"})
		return
	}

	// Return the response as JSON
	c.JSON(http.StatusOK, aircraft)
}

// handleAPIConfigProxy proxies requests to the backend API for configuration
func (s *Server) handleAPIConfigProxy(c *gin.Context) {
	// Create a new request to forward to the API
	req, err := http.NewRequest("GET", s.config.APIEndpoint+"/api/config", nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request for config API"})
		return
	}

	// User is already authenticated via middleware
	// Use the web UI credentials for API auth
	req.SetBasicAuth(s.auth.Username, s.auth.Password)

	// Execute the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch config data from API"})
		return
	}
	defer resp.Body.Close()

	// If the API returns unauthorized, pass that status back to the client
	if resp.StatusCode == http.StatusUnauthorized {
		c.Header("WWW-Authenticate", resp.Header.Get("WWW-Authenticate"))
		c.Status(http.StatusUnauthorized)
		return
	}

	// Parse the response
	var config interface{}
	if err := json.NewDecoder(resp.Body).Decode(&config); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse API config response"})
		return
	}

	// Return the response as JSON
	c.JSON(http.StatusOK, config)
}

// handleVersion serves the application version information
func (s *Server) handleVersion(c *gin.Context) {
	c.JSON(http.StatusOK, version.GetFullVersionInfo())
}

// Start runs the web server
func (s *Server) Start() error {
	addr := fmt.Sprintf(":%d", s.config.ListenPort)
	log.Printf("Starting web frontend at http://localhost%s", addr)

	// Setup template rendering
	tmplFS, err := fs.Sub(content, "templates")
	if err != nil {
		return fmt.Errorf("failed to create template filesystem: %v", err)
	}

	// Create HTML template renderer
	htmlRenderer, err := newTemplateRenderer(tmplFS)
	if err != nil {
		return fmt.Errorf("failed to create template renderer: %v", err)
	}
	s.engine.HTMLRender = htmlRenderer

	// Start the server
	return s.engine.Run(addr)
}

// templateRenderer implements gin.HTMLRender interface
type templateRenderer struct {
	templates map[string]*template.Template
}

// newTemplateRenderer creates a new template renderer with templates from the given filesystem
func newTemplateRenderer(templateFS fs.FS) (render.HTMLRender, error) {
	r := &templateRenderer{
		templates: make(map[string]*template.Template),
	}

	// Load templates from the embedded filesystem
	templateFiles := []string{"index.html", "config.html", "login.html", "loading.html"}
	for _, file := range templateFiles {
		tmpl, err := template.ParseFS(templateFS, file)
		if err != nil {
			return nil, fmt.Errorf("failed to parse template %s: %v", file, err)
		}
		r.templates[file] = tmpl
	}

	return r, nil
}

// Instance implements the gin.HTMLRender interface
func (t *templateRenderer) Instance(name string, data interface{}) render.Render {
	return &templateInstance{
		Name:     name,
		Data:     data,
		Template: t.templates[name],
	}
}

// templateInstance represents a single template instance with data
type templateInstance struct {
	Name     string
	Data     interface{}
	Template *template.Template
}

// Render implements the gin.Render interface
func (t *templateInstance) Render(w http.ResponseWriter) error {
	t.WriteContentType(w)
	return t.Template.Execute(w, t.Data)
}

// WriteContentType writes the content type header to the response
func (t *templateInstance) WriteContentType(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
}
