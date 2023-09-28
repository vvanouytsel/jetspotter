package configuration

import (
	"os"
	"strconv"

	"github.com/jftuga/geodist"
)

type Config struct {
	Location                geodist.Coord
	MaxRangeKilometers      int
	AircraftType            string
	MaxAircraftSlackMessage int
}

// getEnvVariable looks up a specified environment variable, if not set the specified default is used
func getEnvVariable(key, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		value = fallback
	}
	return value
}

func GetConfig() (config Config, err error) {

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
