package web

import (
	"crypto/rand"
	"embed"
	"encoding/json"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"

	"jetspotter/internal/auth"
	"jetspotter/internal/configuration"
	"jetspotter/internal/jetspotter"
	"jetspotter/internal/version"
)

// Templates and static assets embedded in the binary
//
//go:embed templates static
var content embed.FS

// Page names. Each one maps to templates/pages/<name>.html
const (
	pageDashboard = "dashboard"
	pageConfig    = "config"
	pageLogin     = "login"
	pageLoading   = "loading"
)

var allPages = []string{pageDashboard, pageConfig, pageLogin, pageLoading}

// sessionUserKey is the session key that holds the name of the logged in user
const sessionUserKey = "user"

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
	jetspotterConfig *configuration.Config
}

// pageData is everything a template can use. Every page gets the same
// struct so the shared layout (header, footer) always has what it needs.
type pageData struct {
	Page       string // which page is shown, used to highlight the navigation
	Title      string
	IsLoggedIn bool
	Username   string
	Version    string

	// Location of the spotter, used for the "live map" link
	Latitude  float64
	Longitude float64
	MapURL    string

	// Dashboard only
	RefreshSeconds int
	ScanRadius     int

	// Login only
	LoginFailed bool
}

// NewServer creates a new web frontend server
func NewServer(config Config) *Server {
	gin.SetMode(gin.ReleaseMode)

	engine := gin.New()
	engine.Use(gin.Recovery())

	// The session lives in an encrypted and signed cookie, so nothing is kept on the server.
	// A cookie that cannot be verified (for example after a restart) is simply ignored.
	signingKey, encryptionKey := generateSecureKey(), generateSecureKey()
	store := cookie.NewStore(signingKey, encryptionKey)
	store.Options(sessions.Options{
		Path:     "/",
		MaxAge:   86400, // 1 day
		HttpOnly: true,
		Secure:   config.SecureCookies,
		SameSite: http.SameSiteLaxMode,
	})
	engine.Use(sessions.Sessions("jetspotter_session", store))

	jetspotterConfig, err := configuration.GetConfig()
	if err != nil {
		log.Printf("Warning: Failed to load jetspotter configuration: %v", err)
	}

	server := &Server{
		config:           config,
		engine:           engine,
		auth:             auth.NewBasicAuth(),
		jetspotterConfig: &jetspotterConfig,
	}

	engine.Use(requestLogger())
	server.setupRoutes()

	if config.DataReadyChan != nil {
		go server.waitForData(config.DataReadyChan)
	}

	return server
}

// requestLogger logs every request except the ones for static files
func requestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		if strings.HasPrefix(c.Request.URL.Path, "/static/") {
			c.Next()
			return
		}
		start := time.Now()
		c.Next()
		log.Printf("[GIN] %d | %s | %s", c.Writer.Status(), time.Since(start), c.Request.URL.Path)
	}
}

// generateSecureKey creates a random 32 byte key for the session cookie.
// Sessions do not survive a restart of the app, which just means logging in again.
func generateSecureKey() []byte {
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		log.Fatalf("Failed to generate session key: %v", err)
	}
	return key
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
	staticFS, err := fs.Sub(content, "static")
	if err != nil {
		log.Fatalf("Failed to create static file server: %v", err)
	}
	s.engine.StaticFS("/static", http.FS(staticFS))
	s.engine.GET("/favicon.ico", func(c *gin.Context) {
		c.FileFromFS("images/favicon.ico", http.FS(staticFS))
	})

	s.engine.GET("/login", s.handleLoginGet)
	s.engine.POST("/login", s.handleLoginPost)
	s.engine.GET("/logout", s.handleLogout)

	s.engine.GET("/", s.handleDashboard)
	s.engine.GET("/api/aircraft", s.handleAPIProxy)
	s.engine.GET("/api/version", s.handleVersion)

	protected := s.engine.Group("/")
	protected.Use(s.authRequired())
	{
		protected.GET("/config", s.handleConfig)
		protected.GET("/api/config", s.handleAPIConfigProxy)
	}

	s.engine.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"message": "Page not found"})
	})
}

// currentUser returns the name of the logged in user.
// A session that holds something unexpected is cleared.
func currentUser(c *gin.Context) (string, bool) {
	session := sessions.Default(c)
	value := session.Get(sessionUserKey)
	if value == nil {
		return "", false
	}
	username, ok := value.(string)
	if !ok {
		log.Printf("Warning: user session value is not a string type")
		session.Clear()
		session.Save()
		return "", false
	}
	return username, true
}

// newPageData fills in what every page needs
func (s *Server) newPageData(c *gin.Context, page, title string) pageData {
	data := pageData{Page: page, Title: title, Version: displayVersion()}
	data.Username, data.IsLoggedIn = currentUser(c)

	if s.jetspotterConfig != nil {
		data.Latitude = s.jetspotterConfig.Location.Lat
		data.Longitude = s.jetspotterConfig.Location.Lon
	}
	data.MapURL = fmt.Sprintf(
		"https://globe.airplanes.live/?lat=%[1]f&lon=%[2]f&SiteLat=%[1]f&SiteLon=%[2]f&zoom=11&enableLabels&extendedLabels=1&hideSidebar",
		data.Latitude, data.Longitude,
	)
	return data
}

// displayVersion is the version as shown in the footer: v1.2.3, or dev for local builds
func displayVersion() string {
	v := version.GetVersion()
	if v != "" && v[0] >= '0' && v[0] <= '9' {
		return "v" + v
	}
	return v
}

// authRequired is a middleware that redirects to the login page when nobody is logged in
func (s *Server) authRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		if _, loggedIn := currentUser(c); !loggedIn {
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
			return
		}
		c.Next()
	}
}

func (s *Server) handleLoginGet(c *gin.Context) {
	if _, loggedIn := currentUser(c); loggedIn {
		c.Redirect(http.StatusFound, "/")
		return
	}

	data := s.newPageData(c, pageLogin, "Login - Jetspotter")
	data.LoginFailed = c.Query("error") == "true"
	c.HTML(http.StatusOK, pageLogin, data)
}

func (s *Server) handleLoginPost(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	if username != s.auth.Username || password != s.auth.Password {
		c.Redirect(http.StatusFound, "/login?error=true")
		return
	}

	session := sessions.Default(c)
	session.Set(sessionUserKey, username)
	session.Save()
	c.Redirect(http.StatusFound, "/")
}

func (s *Server) handleLogout(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Save()
	c.Redirect(http.StatusFound, "/login")
}

func (s *Server) handleDashboard(c *gin.Context) {
	if !s.isDataReady && s.config.DataReadyChan != nil {
		c.HTML(http.StatusOK, pageLoading, s.newPageData(c, pageLoading, "Jetspotter - Starting"))
		return
	}

	data := s.newPageData(c, pageDashboard, "Jetspotter")
	data.RefreshSeconds = int(s.config.RefreshPeriod.Seconds())
	data.ScanRadius = 30
	if s.jetspotterConfig != nil {
		data.ScanRadius = s.jetspotterConfig.MaxScanRangeKilometers
	}
	c.HTML(http.StatusOK, pageDashboard, data)
}

func (s *Server) handleConfig(c *gin.Context) {
	c.HTML(http.StatusOK, pageConfig, s.newPageData(c, pageConfig, "Jetspotter Configuration"))
}

// handleAPIProxy proxies requests to the backend API
func (s *Server) handleAPIProxy(c *gin.Context) {
	resp, err := http.Get(s.config.APIEndpoint + "/api/aircraft")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch data from API"})
		return
	}
	defer resp.Body.Close()

	var aircraft []jetspotter.Aircraft
	if err := json.NewDecoder(resp.Body).Decode(&aircraft); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse API response"})
		return
	}

	c.JSON(http.StatusOK, aircraft)
}

// handleAPIConfigProxy proxies requests to the backend API for configuration
func (s *Server) handleAPIConfigProxy(c *gin.Context) {
	req, err := http.NewRequest("GET", s.config.APIEndpoint+"/api/config", nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request for config API"})
		return
	}

	// The user is already authenticated by the middleware, so use the web UI credentials for the API
	req.SetBasicAuth(s.auth.Username, s.auth.Password)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch config data from API"})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		c.Header("WWW-Authenticate", resp.Header.Get("WWW-Authenticate"))
		c.Status(http.StatusUnauthorized)
		return
	}

	var config map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&config); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse API config response"})
		return
	}

	maskSecrets(config)
	c.JSON(http.StatusOK, config)
}

// secretConfigKeys are settings that hold webhooks or tokens. The browser only
// needs to know whether they are set, never their value.
var secretConfigKeys = []string{"SlackWebHookURL", "DiscordWebHookURL", "GotifyToken", "NtfyToken"}

func maskSecrets(config map[string]interface{}) {
	for _, key := range secretConfigKeys {
		if value, ok := config[key].(string); ok && value != "" {
			config[key] = "configured"
		}
	}
}

func (s *Server) handleVersion(c *gin.Context) {
	c.JSON(http.StatusOK, version.GetFullVersionInfo())
}

// Start runs the web server
func (s *Server) Start() error {
	addr := fmt.Sprintf(":%d", s.config.ListenPort)
	log.Printf("Starting web frontend at http://localhost%s", addr)

	templateFS, err := fs.Sub(content, "templates")
	if err != nil {
		return fmt.Errorf("failed to create template filesystem: %v", err)
	}

	renderer, err := newTemplateRenderer(templateFS)
	if err != nil {
		return fmt.Errorf("failed to create template renderer: %v", err)
	}
	s.engine.HTMLRender = renderer

	return s.engine.Run(addr)
}

// templateRenderer renders pages. Every page is parsed together with the
// shared layout and partials, so pages only contain what is unique to them.
type templateRenderer struct {
	pages map[string]*template.Template
}

func newTemplateRenderer(templateFS fs.FS) (render.HTMLRender, error) {
	r := &templateRenderer{pages: make(map[string]*template.Template)}

	for _, page := range allPages {
		tmpl, err := template.ParseFS(templateFS, "layouts/*.html", "partials/*.html", "pages/"+page+".html")
		if err != nil {
			return nil, fmt.Errorf("failed to parse page %s: %v", page, err)
		}
		r.pages[page] = tmpl
	}

	return r, nil
}

// Instance implements gin's HTMLRender interface. The name is a page name.
func (r *templateRenderer) Instance(name string, data interface{}) render.Render {
	return &pageRender{template: r.pages[name], data: data}
}

type pageRender struct {
	template *template.Template
	data     interface{}
}

func (p *pageRender) Render(w http.ResponseWriter) error {
	p.WriteContentType(w)
	return p.template.ExecuteTemplate(w, "base", p.data)
}

func (p *pageRender) WriteContentType(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
}
