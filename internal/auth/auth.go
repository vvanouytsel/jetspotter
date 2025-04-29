package auth

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

// BasicAuth represents basic authentication credentials
type BasicAuth struct {
	Username string
	Password string
}

// NewBasicAuth creates a new BasicAuth instance
func NewBasicAuth() *BasicAuth {
	// Get credentials from environment variables or use defaults
	username := "admin"

	password := os.Getenv("AUTH_PASSWORD")
	if password == "" {
		password = "jetspotter" // Default password
	}

	return &BasicAuth{
		Username: username,
		Password: password,
	}
}

// Middleware returns a Gin middleware that implements basic auth
func (a *BasicAuth) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, pass, ok := c.Request.BasicAuth()
		if !ok || user != a.Username || pass != a.Password {
			c.Header("WWW-Authenticate", "Basic realm=Jetspotter API")
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		c.Next()
	}
}

// SetupSessionStore configures the session middleware
func SetupSessionStore() gin.HandlerFunc {
	store := cookie.NewStore([]byte("jetspotter_secret"))
	store.Options(sessions.Options{
		Path:     "/",
		MaxAge:   86400, // 1 day
		HttpOnly: true,
	})
	return sessions.Sessions("jetspotter_session", store)
}

// AuthRequired is a middleware to check if user is logged in
func (a *BasicAuth) AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		user := session.Get("user")
		if user == nil {
			// User is not logged in, redirect to login page
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
			return
		}
		// Continue down the chain
		c.Next()
	}
}

// HandleLogin handles the login form submission
func (a *BasicAuth) HandleLogin(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	// Validate credentials
	if username == a.Username && password == a.Password {
		session := sessions.Default(c)
		session.Set("user", username)
		session.Set("login_time", time.Now().Format(time.RFC3339))
		err := session.Save()
		if err != nil {
			c.HTML(http.StatusInternalServerError, "login.html", gin.H{
				"Title": "Login - Jetspotter",
				"Error": "Failed to create session",
			})
			return
		}
		// Redirect to home page after successful login
		c.Redirect(http.StatusFound, "/")
		return
	}

	// Invalid credentials
	c.Redirect(http.StatusFound, "/login?error=true")
}

// HandleLogout logs the user out
func (a *BasicAuth) HandleLogout(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Save()
	c.Redirect(http.StatusFound, "/login")
}

// GetCurrentUser returns the currently logged in username
func (a *BasicAuth) GetCurrentUser(c *gin.Context) string {
	session := sessions.Default(c)
	user := session.Get("user")
	if user == nil {
		return ""
	}
	return user.(string)
}

// IsUserLoggedIn checks if a user is currently logged in
func (a *BasicAuth) IsUserLoggedIn(c *gin.Context) bool {
	session := sessions.Default(c)
	user := session.Get("user")
	return user != nil
}
