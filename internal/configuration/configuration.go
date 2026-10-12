package configuration

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/jftuga/geodist"
)

// Config is the type used for all user configuration.
// All parameters can be set using ENV variables.
// The comments below are structured as following:
// ENV_VARIABLE_NAME DEFAULT_VALUE
type Config struct {
	// Latitude and Longitude coordinates of the location you want to use.
	// LOCATION_LATITUDE 51.17348
	// LOCATION_LONGITUDE 5.45921
	Location geodist.Coord

	// Maximum range in kilometers from the location that you want aircraft to be spotted.
	// Note that this is an approximation due to roundings.
	// MAX_RANGE_KILOMETERS 30
	MaxRangeKilometers int

	// Maximum range in kilometers from the location that you want to scan for aircraft.
	// This is used to determine the area to query from the ADS-B API.
	// If not set, MAX_RANGE_KILOMETERS will be used.
	// MAX_SCAN_RANGE_KILOMETERS 30
	MaxScanRangeKilometers int

	// Maximum altitude in feet that you want to spot aircraft at.
	// Set to 0 to disable the filter.
	// MAX_ALTITUDE_FEET 0
	MaxAltitudeFeet int

	// A comma seperated list of types that you want to spot
	// If not set, 'ALL' will be used, which will disable the filter and show all aircraft within range.
	// Full list can be found at https://www.icao.int/publications/doc8643/pages/search.aspx in 'Type Designator' column.
	// AIRCRAFT_TYPES ALL
	// EXAMPLES
	// AIRCRAFT_TYPES F16,F35
	// To spot all military aircraft, you can use MILITARY.
	// AIRCRAFT_TYPES MILITARY
	AircraftTypes []string

	// Webhook used to send notifications to Slack. If not set, no messages will be sent to Slack.
	// SLACK_WEBHOOK_URL ""
	SlackWebHookURL string

	// Webhook used to send notifications to Discord. If not set, no messages will be sent to Discord.
	// DISCORD_WEBHOOK_URL ""
	DiscordWebHookURL string

	// Discord notifications use an embed color based on the alitute of the aircraft.
	// DISCORD_COLOR_ALTITUDE "true"
	DiscordColorAltitude string

	// Interval in seconds between fetching aircraft, minimum is 60 due to API rate limiting.
	// FETCH_INTERVAL 60
	FetchInterval int

	// Token to authenticate with the gotify server.
	// GOTIFY_TOKEN ""
	GotifyToken string

	// URL of the gotify server.
	// GOTIFY_URL ""
	GotifyURL string

	// Port where metrics will be exposed on
	// METRICS_PORT "7070"
	MetricsPort string

	// Port where API will be exposed on
	// API_PORT "8085"
	APIPort string

	// Enable or disable the web UI
	// WEB_UI_ENABLED "true"
	WebUIEnabled bool

	// Port where web UI will be exposed on
	// WEB_UI_PORT "8080"
	WebUIPort string

	// Topic to publish message to
	// NTFY_TOPIC ""
	NtfyTopic string

	// URL of the ntfy server.
	// NTFY_SERVER "https://ntfy.sh"
	NtfyServer string

	// Token for ntfy server authentication.
	// NTFY_TOKEN ""
	NtfyToken string

	// Enable overhead trajectory prediction. When enabled, the application
	// predicts aircraft that are on course to fly directly overhead and sends
	// a "look up in N minutes" notification after a confirmation window.
	// OVERHEAD_PREDICTION_ENABLED "false"
	OverheadPredictionEnabled bool

	// Lateral radius in kilometers from the location that defines "overhead".
	// An aircraft whose predicted closest point of approach is within this
	// radius (and that stays on course for the confirmation window) triggers
	// a predictive notification.
	// OVERHEAD_RADIUS_KILOMETERS 3
	OverheadRadiusKilometers int

	// Look-ahead horizon in minutes. Only aircraft whose predicted closest
	// point of approach occurs within this many minutes from now are
	// considered candidates. Larger values give more warning but produce more
	// (eventually abandoned) candidates.
	// OVERHEAD_LOOK_AHEAD_MINUTES 10
	OverheadLookAheadMinutes int

	// Confirmation window in seconds. A candidate must remain on course and
	// projected overhead for this many seconds of consecutive sub-loop samples
	// before a notification is sent.
	// OVERHEAD_CONFIRM_SECONDS 5
	OverheadConfirmSeconds int

	// Polling interval in seconds for the per-candidate confirmation sub-loop.
	// Each candidate is polled at this cadence until confirmed or abandoned.
	// OVERHEAD_SUBLOOP_POLL_SECONDS 5
	OverheadSubloopPollSeconds int

	// Inbound margin in degrees. An aircraft is considered "on course" for
	// overhead prediction when its heading differs from the bearing to the
	// target location by at most this many degrees. Smaller is stricter.
	// OVERHEAD_INBOUND_MARGIN_DEGREES 10
	OverheadInboundMarginDegrees int
}

// Environment variable names
const (
	SlackWebhookURL        = "SLACK_WEBHOOK_URL"
	DiscordWebhookURL      = "DISCORD_WEBHOOK_URL"
	DiscordColorAltitude   = "DISCORD_COLOR_ALTITUDE"
	LocationLatitude       = "LOCATION_LATITUDE"
	LocationLongitude      = "LOCATION_LONGITUDE"
	MaxRangeKilometers     = "MAX_RANGE_KILOMETERS"
	MaxScanRangeKilometers = "MAX_SCAN_RANGE_KILOMETERS"
	MaxAltitudeFeet        = "MAX_ALTITUDE_FEET"
	AircraftTypes          = "AIRCRAFT_TYPES"
	FetchInterval          = "FETCH_INTERVAL"
	GotifyURL              = "GOTIFY_URL"
	NtfyTopic              = "NTFY_TOPIC"
	NtfyServer             = "NTFY_SERVER"
	NtfyToken              = "NTFY_TOKEN"
	GotifyToken            = "GOTIFY_TOKEN"
	MetricsPort            = "METRICS_PORT"
	APIPort                = "API_PORT"
	WebUIEnabled           = "WEB_UI_ENABLED"
	WebUIPort              = "WEB_UI_PORT"
	OverheadPredictionEnabled  = "OVERHEAD_PREDICTION_ENABLED"
	OverheadRadiusKilometers   = "OVERHEAD_RADIUS_KILOMETERS"
	OverheadLookAheadMinutes   = "OVERHEAD_LOOK_AHEAD_MINUTES"
	OverheadConfirmSeconds    = "OVERHEAD_CONFIRM_SECONDS"
	OverheadSubloopPollSeconds = "OVERHEAD_SUBLOOP_POLL_SECONDS"
	OverheadInboundMarginDeg   = "OVERHEAD_INBOUND_MARGIN_DEGREES"
)

// getEnvVariable looks up a specified environment variable, if not set the specified default is used
func getEnvVariable(key, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		value = fallback
	}
	return value
}

// GetConfig attempts to read the configuration via environment variables and uses a default if the environment variable is not set
func GetConfig() (config Config, err error) {
	defaultFetchInterval := 60

	config.GotifyToken = getEnvVariable(GotifyToken, "")
	config.GotifyURL = getEnvVariable(GotifyURL, "")
	config.NtfyTopic = getEnvVariable(NtfyTopic, "")
	config.NtfyServer = getEnvVariable(NtfyServer, "https://ntfy.sh")
	config.NtfyToken = getEnvVariable(NtfyToken, "")
	config.SlackWebHookURL = getEnvVariable(SlackWebhookURL, "")
	config.DiscordWebHookURL = getEnvVariable(DiscordWebhookURL, "")
	config.DiscordColorAltitude = getEnvVariable(DiscordColorAltitude, "true")
	config.MetricsPort = getEnvVariable(MetricsPort, "7070")
	config.APIPort = getEnvVariable(APIPort, "8085")
	config.WebUIPort = getEnvVariable(WebUIPort, "8080")

	// Convert WebUIEnabled string to bool
	webUIEnabledStr := getEnvVariable(WebUIEnabled, "true")
	config.WebUIEnabled, err = strconv.ParseBool(webUIEnabledStr)
	if err != nil {
		log.Printf("Invalid value for WEB_UI_ENABLED: %s, using default: true", webUIEnabledStr)
		config.WebUIEnabled = true
	}

	config.FetchInterval, err = strconv.Atoi(getEnvVariable(FetchInterval, strconv.Itoa(defaultFetchInterval)))
	if err != nil || config.FetchInterval < 5 {
		log.Printf("Fetch interval of %ds detected. You might hit rate limits, consider using the default of %ds instead.", config.FetchInterval, defaultFetchInterval)
	}

	config.Location.Lat, err = strconv.ParseFloat(getEnvVariable(LocationLatitude, "51.17348"), 64)
	if err != nil {
		return Config{}, err
	}

	config.Location.Lon, err = strconv.ParseFloat(getEnvVariable(LocationLongitude, "5.45921"), 64)
	if err != nil {
		return Config{}, err
	}

	config.MaxRangeKilometers, err = strconv.Atoi(getEnvVariable(MaxRangeKilometers, "30"))
	if err != nil {
		return Config{}, err
	}

	// Get MAX_SCAN_RANGE_KILOMETERS - If not set, use MAX_RANGE_KILOMETERS for backward compatibility
	maxScanRangeStr := getEnvVariable(MaxScanRangeKilometers, "")
	if maxScanRangeStr == "" {
		config.MaxScanRangeKilometers = config.MaxRangeKilometers
	} else {
		config.MaxScanRangeKilometers, err = strconv.Atoi(maxScanRangeStr)
		if err != nil {
			return Config{}, err
		}
	}

	config.MaxAltitudeFeet, err = strconv.Atoi(getEnvVariable(MaxAltitudeFeet, "0"))
	if err != nil {
		return Config{}, err
	}

	config.AircraftTypes = strings.Split(strings.ToUpper(strings.ReplaceAll(getEnvVariable(AircraftTypes, "ALL"), " ", "")), ",")

	config.OverheadPredictionEnabled, err = strconv.ParseBool(getEnvVariable(OverheadPredictionEnabled, "false"))
	if err != nil {
		log.Printf("Invalid value for OVERHEAD_PREDICTION_ENABLED: %s, using default: false", getEnvVariable(OverheadPredictionEnabled, "false"))
		config.OverheadPredictionEnabled = false
	}

	config.OverheadRadiusKilometers, err = strconv.Atoi(getEnvVariable(OverheadRadiusKilometers, "3"))
	if err != nil {
		log.Printf("Invalid value for OVERHEAD_RADIUS_KILOMETERS: %s, using default: 3", getEnvVariable(OverheadRadiusKilometers, "3"))
		config.OverheadRadiusKilometers = 3
	}

	config.OverheadLookAheadMinutes, err = strconv.Atoi(getEnvVariable(OverheadLookAheadMinutes, "10"))
	if err != nil {
		log.Printf("Invalid value for OVERHEAD_LOOK_AHEAD_MINUTES: %s, using default: 10", getEnvVariable(OverheadLookAheadMinutes, "10"))
		config.OverheadLookAheadMinutes = 10
	}

	config.OverheadConfirmSeconds, err = strconv.Atoi(getEnvVariable(OverheadConfirmSeconds, "5"))
	if err != nil || config.OverheadConfirmSeconds < 1 {
		log.Printf("Invalid value for OVERHEAD_CONFIRM_SECONDS: %s, using default: 5", getEnvVariable(OverheadConfirmSeconds, "5"))
		config.OverheadConfirmSeconds = 5
	}

	config.OverheadSubloopPollSeconds, err = strconv.Atoi(getEnvVariable(OverheadSubloopPollSeconds, "5"))
	if err != nil || config.OverheadSubloopPollSeconds < 1 {
		log.Printf("Invalid value for OVERHEAD_SUBLOOP_POLL_SECONDS: %s, using default: 5", getEnvVariable(OverheadSubloopPollSeconds, "5"))
		config.OverheadSubloopPollSeconds = 5
	}

	config.OverheadInboundMarginDegrees, err = strconv.Atoi(getEnvVariable(OverheadInboundMarginDeg, "10"))
	if err != nil || config.OverheadInboundMarginDegrees < 1 || config.OverheadInboundMarginDegrees > 180 {
		log.Printf("Invalid value for OVERHEAD_INBOUND_MARGIN_DEGREES: %s, using default: 10", getEnvVariable(OverheadInboundMarginDeg, "10"))
		config.OverheadInboundMarginDegrees = 10
	}

	return config, nil
}
