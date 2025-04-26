package version

// These variables will be populated at build time
var (
	// Version holds the current version of the application
	Version = "dev"
	
	// Commit holds the git commit hash used to build the application
	Commit = "none"
	
	// BuildTime holds the time when the application was built
	BuildTime = "unknown"
)

// GetVersion returns the full version information as a string
func GetVersion() string {
	return Version
}

// GetFullVersionInfo returns the complete version information
func GetFullVersionInfo() map[string]string {
	return map[string]string{
		"version":   Version,
		"commit":    Commit,
		"buildTime": BuildTime,
	}
}