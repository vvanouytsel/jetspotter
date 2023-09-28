package jetspotter

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"

	aircr "jetspotter/internal/aircraft"
	"jetspotter/internal/configuration"

	"github.com/jftuga/geodist"
)

// Vars
var (
	baseURL = "https://api.adsb.one/v2/point"
)

/* CalculateDistance returns the rounded distance between two coordinates in kilometers */
func CalculateDistance(source geodist.Coord, destination geodist.Coord) int {
	_, kilometers := geodist.HaversineDistance(source, destination)
	return int(kilometers)
}

/* GetAircraftInProximity returns all aircraft within a specified maxRange of a latitude/longitude point. */
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

// GetFiltererdAircraftInRange returns all aircraft of specified type within maxRange kilometers of the location.
func GetFiltererdAircraftInRange(location geodist.Coord, aircraftType string, maxRange int) (aircraft []Aircraft, err error) {
	var flightData FlightData
	miles := int(float32(maxRange) / 1.60934)
	endpoint, err := url.JoinPath(baseURL,
		strconv.FormatFloat(location.Lat, 'f', -1, 64),
		strconv.FormatFloat(location.Lon, 'f', -1, 64),
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

	if aircraftType == aircr.ALL.Identifier {
		return flightData.AC, nil
	}

	return filterAircraftByType(flightData.AC, aircraftType), nil
}

// filterAircraftByType returns a list of Aircraft that match the aircraftType.
func filterAircraftByType(aircraft []Aircraft, aircraftType string) []Aircraft {
	var filteredAircraft []Aircraft

	for _, ac := range aircraft {
		if ac.PlaneType == aircraftType || aircraftType == aircr.ALL.Identifier {
			filteredAircraft = append(filteredAircraft, ac)
		}
	}
	return filteredAircraft
}

// FormatAircraft prints an Aircraft in a readable manner.
func FormatAircraft(aircraft Aircraft, config configuration.Config) string {
	if aircraft.Callsign == "" {
		aircraft.Callsign = "UNKNOWN"
	}

	distance := CalculateDistance(
		config.Location,
		geodist.Coord{
			Lat: aircraft.Lat,
			Lon: aircraft.Lon,
		},
	)

	return fmt.Sprintf("Callsign: %s\n"+
		"Description: %s\n"+
		"Type: %s\n"+
		"Tail number: %s\n"+
		"Altitude: %vft\n"+
		"Speed: %dkn\n"+
		"Distance: %vkm\n"+
		"URL: %s",
		aircraft.Callsign, aircraft.Desc, aircraft.PlaneType,
		aircraft.TailNumber, aircraft.AltBaro,
		int(aircraft.GS), distance, fmt.Sprintf("https://globe.adsbexchange.com/?icao=%s\n", aircraft.ICAO))
}

// PrintAircraft prints a list of Aircraft in a readable manner.
func PrintAircraft(aircraft []Aircraft, config configuration.Config) {
	if len(aircraft) == 0 {
		fmt.Println("No matching aircraft have been spotted.")
	}

	for _, ac := range aircraft {
		fmt.Println(FormatAircraft(ac, config))
	}
}
