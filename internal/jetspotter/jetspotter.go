package jetspotter

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"slices"
	"strconv"
	"time"

	"jetspotter/internal/configuration"
	"jetspotter/internal/weather"

	"github.com/jftuga/geodist"
)

// Vars
var (
	baseURL = "https://api.adsb.one/v2/point"
)

// CalculateDistance returns the rounded distance between two coordinates in kilometers
func CalculateDistance(source geodist.Coord, destination geodist.Coord) int {
	_, kilometers := geodist.HaversineDistance(source, destination)
	return int(kilometers)
}

// GetAircraftInProximity returns all aircraft within a specified maxRange of a latitude/longitude point
func GetAircraftInProximity(latitude string, longitude string, maxRange int) (aircraft []Aircraft, err error) {
	var flightData FlightData
	endpoint, err := url.JoinPath(baseURL, latitude, longitude, strconv.Itoa(maxRange))
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)
	err = json.Unmarshal(body, &flightData)
	if err != nil {
		return nil, err
	}

	return flightData.AC, nil
}

// GetFiltererdAircraftInRange returns all aircraft of specified types within maxRange kilometers of the location.
func GetFiltererdAircraftInRange(config configuration.Config) (aircraft []AircraftOutput, err error) {
	var flightData FlightData
	miles := int(float32(config.MaxRangeKilometers) / 1.60934)
	endpoint, err := url.JoinPath(baseURL,
		strconv.FormatFloat(config.Location.Lat, 'f', -1, 64),
		strconv.FormatFloat(config.Location.Lon, 'f', -1, 64),
		strconv.Itoa(miles))
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)
	err = json.Unmarshal(body, &flightData)
	if err != nil {
		return nil, err
	}

	acOutputs, err := CreateAircraftOutput(flightData.AC, config)
	if err != nil {
		return nil, err
	}

	if slices.Contains(config.AircraftTypes, "ALL") {
		return acOutputs, nil
	}
	return filterAircraftByTypes(acOutputs, config), nil
}

// filterAircraftByTypes returns a list of Aircraft that match the aircraftTypes.
func filterAircraftByTypes(aircraft []AircraftOutput, config configuration.Config) []AircraftOutput {
	var filteredAircraft []AircraftOutput

	for _, ac := range aircraft {
		for _, aircraftType := range config.AircraftTypes {
			if ac.Type == aircraftType || aircraftType == "ALL" {
				filteredAircraft = append(filteredAircraft, ac)
			}
		}
	}

	return filteredAircraft
}

// FormatAircraft prints an Aircraft in a readable manner.
func FormatAircraft(aircraft AircraftOutput, config configuration.Config) string {

	return fmt.Sprintf("Callsign: %s\n"+
		"Description: %s\n"+
		"Type: %s\n"+
		"Tail number: %s\n"+
		"Altitude: %dft | %dm\n"+
		"Speed: %dkn | %dkm/h\n"+
		"Distance: %dkm\n"+
		"Cloud coverage: %d%%\n"+
		"URL: %s",
		aircraft.Callsign, aircraft.Description, aircraft.Type,
		aircraft.TailNumber, int(aircraft.Altitude), ConvertFeetToMeters(aircraft.Altitude),
		aircraft.Speed, ConvertKnotsToKilometersPerHour(aircraft.Speed), aircraft.Distance, aircraft.CloudCoverage, aircraft.URL)
}

// PrintAircraft prints a list of Aircraft in a readable manner.
func PrintAircraft(aircraft []AircraftOutput, config configuration.Config) {
	if len(aircraft) == 0 {
		fmt.Println("No matching aircraft have been spotted.")
	}

	for _, ac := range aircraft {
		fmt.Println(FormatAircraft(ac, config))
	}
}

// ConvertKnotsToKilometersPerHour well converts knots to kilometers per hour...
func ConvertKnotsToKilometersPerHour(knots int) int {
	return int(float64(knots) * 1.852)
}

// ConvertFeetToMeters converts feet to meters, * pikachu face *
func ConvertFeetToMeters(feet float64) int {
	return int(feet * 0.3048)
}

// getCloudCoverage gets the coverage percentage of the clouds at a given altitude block
// Altitude blocks are one of the following
// low    -> 0m up to 3000m
// medium -> 3000m up to 8000m
// high   -> above 8000m
func getCloudCoverage(weather weather.WeatherData, altitudeInFeet float64) (cloudCoveragePercentage int) {

	altitudeInMeters := ConvertFeetToMeters(altitudeInFeet)
	hourUTC := (time.Now().Hour())

	switch {
	case altitudeInMeters < 3000:
		return weather.Hourly.CloudcoverLow[hourUTC]
	case altitudeInMeters >= 3000 && altitudeInMeters < 8000:
		return weather.Hourly.CloudcoverMid[hourUTC]
	default:
		return weather.Hourly.CloudcoverHigh[hourUTC]
	}
}

func validateFields(aircraft Aircraft) Aircraft {
	if aircraft.Callsign == "" {
		aircraft.Callsign = "UNKNOWN"
	}

	if aircraft.AltBaro == "groundft" || aircraft.AltBaro == "ground" {
		aircraft.AltBaro = 0
	}

	altitudeBarometricFloat := aircraft.AltBaro.(float64)
	if altitudeBarometricFloat < 0 {
		altitudeBarometricFloat = 0
		aircraft.AltBaro = altitudeBarometricFloat
	}

	return aircraft
}

// CreateAircraftOutput returns a list of AircraftOutput objects that will be used to print metadata.
func CreateAircraftOutput(aircraft []Aircraft, config configuration.Config) (acOutputs []AircraftOutput, err error) {
	var acOutput AircraftOutput

	weather, err := weather.GetCloudForecast(config.Location)
	if err != nil {
		return nil, err
	}

	// Wrap these checks in a function
	for _, ac := range aircraft {
		ac = validateFields(ac)

		acOutput.Altitude = ac.AltBaro.(float64)
		acOutput.Callsign = ac.Callsign
		acOutput.Description = ac.Desc
		acOutput.Distance = CalculateDistance(
			config.Location,
			geodist.Coord{
				Lat: ac.Lat,
				Lon: ac.Lon,
			},
		)
		acOutput.Speed = int(ac.GS)
		acOutput.TailNumber = ac.TailNumber
		acOutput.Type = ac.PlaneType
		acOutput.URL = fmt.Sprintf("https://globe.adsbexchange.com/?icao=%s\n", ac.ICAO)
		acOutput.CloudCoverage = getCloudCoverage(*weather, acOutput.Altitude)

		acOutputs = append(acOutputs, acOutput)
	}
	return acOutputs, nil
}
