package configuration

import (
	"os"
	"strconv"

	"github.com/jftuga/geodist"
)

// Config is the type used for all user configuration.
// All parameters can be set using ENV variables.
// The comments below are structured as following:
// ENV_VARIABLE_NAME DEFAULT_VALUE
type Config struct {
	// Latitude and Longitude coordinates of the location you want to use.
	// LOCATION_LATITUDE 51.078395
	// LOCATION_LONGITUDE 5.018769
	Location geodist.Coord

	// Maximum range in kilometers from the location that you want aircraft to be spotted.
	// Note that this is an approximation due to roundings.
	// MAX_RANGE_KILOMETER 30
	MaxRangeKilometers int

	// Type of aircraft that you want to spot.
	// A list of types is not yet supported, however using 'ALL' will disable the filter and show all aircraft.
	// Full list can be found at https://www.icao.int/publications/doc8643/pages/search.aspx in 'Type Designator' column.
	// AIRCRAFT_TYPE ALL
	AircraftType string

	// Maximum amount of aircraft to show in a single slack message.
	// Note that a single slack message only supports up to 50 'blocks' and each aircraft that we display has multiple blocks.
	// MAX_AIRCRAFT_SLACK_MESSAGE 8
	MaxAircraftSlackMessage int

	// Webhook used to send notifications to Slack. If not set, no messages will be sent to Slack.
	// SLACK_WEBHOOK_URL ""
	SlackWebHookURL string
}

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

	config.SlackWebHookURL = getEnvVariable("SLACK_WEBHOOK_URL", "")

	config.Location.Lat, err = strconv.ParseFloat(getEnvVariable("LOCATION_LATITUDE", "51.078395"), 64)
	if err != nil {
		return Config{}, err
	}

	config.Location.Lon, err = strconv.ParseFloat(getEnvVariable("LOCATION_LONGITUDE", "5.018769"), 64)
	if err != nil {
		return Config{}, err
	}

	config.MaxRangeKilometers, err = strconv.Atoi(getEnvVariable("MAX_RANGE_KILOMETERS", "30"))
	if err != nil {
		return Config{}, err
	}

	config.MaxAircraftSlackMessage, err = strconv.Atoi(getEnvVariable("MAX_AIRCRAFT_SLACK_MESSAGE", "8"))
	if err != nil {
		return Config{}, err
	}

	config.AircraftType = getEnvVariable("AIRCRAFT_TYPE", "ALL")
	return config, nil

}
