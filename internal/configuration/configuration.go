package configuration

import (
	"fmt"
	"log"
	"strings"

	"github.com/jftuga/geodist"
	"github.com/spf13/viper"
)

var (
	AircraftTypes           []string
	DiscordWebhookURL       string
	DiscordColorAltitude    bool
	FetchInterval           int
	GotifyToken             string
	GotifyURL               string
	Location                geodist.Coord
	LogNewPlanesToConsole   bool
	MaxAircraftSlackMessage int
	MaxRangeKilometers      int
	MetricsPort             int
	NtfyServer              string
	NtfyTopic               string
	SlackWebhookURL         string
)

const (
	// Defaults

	defaultFetchInterval           int = 60
	defaultMaxAircraftSlackMessage int = 8
	defaultMaxRangeKilometers      int = 30

	// Config IDs

	config_aircraftTypes           string = "aircraftTypes"
	config_discordColorAltitude    string = "discordColorAltitude"
	config_discordWebHookUrl       string = "discordWebHookUrl"
	config_fetchInterval           string = "fetchInterval"
	config_gotifyToken             string = "gotifyToken"
	config_gotifyUrl               string = "gotifyUrl"
	config_locationLatitude        string = "locationLatitude"
	config_locationLongitude       string = "locationLongitude"
	config_logNewPlanesToConsole   string = "logNewPlanesToConsole"
	config_maxAircraftSlackMessage string = "maxAircraftSlackMessage"
	config_maxRangeKilometers      string = "maxRangeKilometers"
	config_metricsPort             string = "metricsPort"
	config_ntfyServer              string = "ntfyServer"
	config_ntfyTopic               string = "ntfyTopic"
	config_slackWebhookUrl         string = "slackWebhookUrl"

	// Environment variables

	env_aircraftTypes           string = "AIRCRAFT_TYPES"
	env_discordColorAltitude    string = "DISCORD_COLOR_ALTITUDE"
	env_discordWebHookUrl       string = "DISCORD_WEBHOOK_URL"
	env_fetchInterval           string = "FETCH_INTERVAL"
	env_gotifyToken             string = "GOTIFY_TOKEN"
	env_gotifyUrl               string = "GOTIFY_URL"
	env_locationLatitude        string = "LOCATION_LATITUDE"
	env_locationLongitude       string = "LOCATION_LONGITUDE"
	env_logNewPlanesToConsole   string = "LOG_NEW_PLANES_TO_CONSOLE"
	env_maxAircraftSlackMessage string = "MAX_AIRCRAFT_SLACK_MESSAGE"
	env_maxRangeKilometers      string = "MAX_RANGE_KILOMETERS"
	env_metricsPort             string = "METRICS_PORT"
	env_ntfyServer              string = "NTFY_SERVER"
	env_ntfyTopic               string = "NTFY_TOPIC"
	env_slackWebhookUrl         string = "SLACK_WEBHOOK_URL"
)

// GetConfig attempts to read the configuration via environment variables and uses a default if the environment variable is not set
func GetConfig() {
	var err error

	viper.SetDefault(config_aircraftTypes, "ALL")
	viper.SetDefault(config_discordColorAltitude, true)
	viper.SetDefault(config_discordWebHookUrl, "")
	viper.SetDefault(config_fetchInterval, defaultFetchInterval)
	viper.SetDefault(config_gotifyToken, "")
	viper.SetDefault(config_gotifyUrl, "")
	viper.SetDefault(config_locationLatitude, 51.17348)
	viper.SetDefault(config_locationLongitude, 5.45921)
	viper.SetDefault(config_logNewPlanesToConsole, true)
	viper.SetDefault(config_maxAircraftSlackMessage, defaultMaxAircraftSlackMessage)
	viper.SetDefault(config_maxRangeKilometers, defaultMaxRangeKilometers)
	viper.SetDefault(config_metricsPort, 7070)
	viper.SetDefault(config_ntfyServer, "https://ntfy.sh")
	viper.SetDefault(config_ntfyTopic, "")
	viper.SetDefault(config_slackWebhookUrl, "")

	// Bind the Viper key to an associated environment variable name

	err = viper.BindEnv(config_aircraftTypes, env_aircraftTypes)
	if err != nil {
		log.Fatal(err)
	}

	err = viper.BindEnv(config_discordColorAltitude, env_discordColorAltitude)
	if err != nil {
		log.Fatal(err)
	}

	err = viper.BindEnv(config_discordWebHookUrl, env_discordWebHookUrl)
	if err != nil {
		log.Fatal(err)
	}

	err = viper.BindEnv(config_fetchInterval, env_fetchInterval)
	if err != nil {
		log.Fatal(err)
	}

	err = viper.BindEnv(config_gotifyToken, env_gotifyToken)
	if err != nil {
		log.Fatal(err)
	}

	err = viper.BindEnv(config_gotifyUrl, env_gotifyUrl)
	if err != nil {
		log.Fatal(err)
	}

	err = viper.BindEnv(config_locationLatitude, env_locationLatitude)
	if err != nil {
		log.Fatal(err)
	}

	err = viper.BindEnv(config_locationLongitude, env_locationLongitude)
	if err != nil {
		log.Fatal(err)
	}

	err = viper.BindEnv(config_logNewPlanesToConsole, env_logNewPlanesToConsole)
	if err != nil {
		log.Fatal(err)
	}

	err = viper.BindEnv(config_maxAircraftSlackMessage, env_maxAircraftSlackMessage)
	if err != nil {
		log.Fatal(err)
	}

	err = viper.BindEnv(config_maxRangeKilometers, env_maxRangeKilometers)
	if err != nil {
		log.Fatal(err)
	}

	err = viper.BindEnv(config_metricsPort, env_metricsPort)
	if err != nil {
		log.Fatal(err)
	}

	err = viper.BindEnv(config_ntfyServer, env_ntfyServer)
	if err != nil {
		log.Fatal(err)
	}

	err = viper.BindEnv(config_ntfyTopic, env_ntfyTopic)
	if err != nil {
		log.Fatal(err)
	}

	err = viper.BindEnv(config_slackWebhookUrl, env_slackWebhookUrl)
	if err != nil {
		log.Fatal(err)
	}

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error
		} else {
			panic(fmt.Errorf("fatal error in config file: %w", err))
		}
	}

	aircraftTypes := viper.GetString(config_aircraftTypes)
	AircraftTypes = strings.Split(strings.ToUpper(strings.ReplaceAll(aircraftTypes, " ", "")), ",")

	DiscordColorAltitude = viper.GetBool(config_discordColorAltitude)
	DiscordWebhookURL = viper.GetString(config_discordWebHookUrl)
	FetchInterval = viper.GetInt(config_fetchInterval)
	GotifyToken = viper.GetString(config_gotifyToken)
	GotifyURL = viper.GetString(config_gotifyUrl)
	Location.Lat = viper.GetFloat64(config_locationLatitude)
	Location.Lon = viper.GetFloat64(config_locationLongitude)
	LogNewPlanesToConsole = viper.GetBool(config_logNewPlanesToConsole)
	MaxAircraftSlackMessage = viper.GetInt(config_maxAircraftSlackMessage)
	MaxRangeKilometers = viper.GetInt(config_maxRangeKilometers)
	MetricsPort = viper.GetInt(config_metricsPort)
	NtfyServer = viper.GetString(config_ntfyServer)
	NtfyTopic = viper.GetString(config_ntfyTopic)
	SlackWebhookURL = viper.GetString(config_slackWebhookUrl)

	if FetchInterval < 60 {
		log.Printf("Fetch interval of %ds detected. You might hit rate limits. Please consider using the default of %ds instead.", FetchInterval, defaultFetchInterval)
	}
}
