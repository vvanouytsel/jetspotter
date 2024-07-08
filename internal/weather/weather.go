package weather

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/jftuga/geodist"
)

// HourlyUnits represents the cloud cover units per hour
type HourlyUnits struct {
	Time           string `json:"time"`
	CloudcoverLow  string `json:"cloudcover_low"`
	CloudcoverMid  string `json:"cloudcover_mid"`
	CloudcoverHigh string `json:"cloudcover_high"`
}

// HourlyData represents the cloud cover data per hour
type HourlyData struct {
	Time           []string `json:"time"`
	CloudcoverLow  []int    `json:"cloudcover_low"`
	CloudcoverMid  []int    `json:"cloudcover_mid"`
	CloudcoverHigh []int    `json:"cloudcover_high"`
}

// Data represents the data of the weather
type Data struct {
	Latitude         float64     `json:"latitude"`
	Longitude        float64     `json:"longitude"`
	GenerationTimeMs float64     `json:"generationtime_ms"`
	UtcOffsetSeconds int         `json:"utc_offset_seconds"`
	Timezone         string      `json:"timezone"`
	TimezoneAbbrev   string      `json:"timezone_abbreviation"`
	Elevation        float64     `json:"elevation"`
	HourlyUnits      HourlyUnits `json:"hourly_units"`
	Hourly           HourlyData  `json:"hourly"`
}

const (
	weatherBaseURL = "https://api.open-meteo.com/v1/forecast?"
)

// GetCloudForecast gets the cloud forecast for every hour of the current day based on a specified location
func GetCloudForecast(location geodist.Coord) (weather *Data, err error) {

	weatherCloudURL := weatherBaseURL + fmt.Sprintf("latitude=%.6f&longitude=%.6f&hourly=cloudcover_low,cloudcover_mid,cloudcover_high&windspeed_unit=kn&timezone=GMT&forecast_days=1",
		location.Lat, location.Lon)

	request, err := http.NewRequest(http.MethodGet, weatherCloudURL, nil)
	if err != nil {
		return nil, err
	}

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != 200 {
		return nil, fmt.Errorf("Received status code %v", response.StatusCode)
	}

	err = json.Unmarshal(body, &weather)
	if err != nil {
		return nil, err
	}

	return weather, nil
}
