package jetspotter

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"

	aircr "jetspotter/internal/aircraft"

	"github.com/jftuga/geodist"
)

// Vars
var (
	baseURL  = "https://api.adsb.one/v2/point"
	BullsEye = Location{
		Latitude:  51.078395,
		Longitude: 5.018769,
	}
)

/* calculateDistance returns the rounded distance between two coordinates in kilometers */
func CalculateDistance(source geodist.Coord, destination geodist.Coord) int {
	_, kilometers := geodist.HaversineDistance(source, destination)
	return int(kilometers)
}

/* getAircraftInProximity returns all aircraft within a specified maxRange of a latitude/longitude point. */
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

// GetAircraftTypeInRange returns all aircraft of specified type within maxRange kilometers of the location.
func GetAircraftTypeInRange(location Location, aircraftType string, maxRange int) (aircraft []Aircraft, err error) {
	var flightData FlightData
	miles := int(float32(maxRange) / 1.60934)
	endpoint, err := url.JoinPath(baseURL,
		strconv.FormatFloat(location.Latitude, 'f', -1, 64),
		strconv.FormatFloat(location.Longitude, 'f', -1, 64),
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
func FormatAircraft(aircraft Aircraft) string {
	if aircraft.Flight == "" {
		aircraft.Flight = "UNKNOWN"
	}

	distance := CalculateDistance(
		geodist.Coord{
			Lat: BullsEye.Latitude,
			Lon: BullsEye.Longitude,
		},
		geodist.Coord{
			Lat: aircraft.Lat,
			Lon: aircraft.Lon,
		},
	)

	return fmt.Sprintf("Callsign: %s\nDescription: %s\nType: %s\nTail number: %s\nAltitude: %v\nDistance: %vkm\n",
		aircraft.Flight, aircraft.Desc, aircraft.PlaneType, aircraft.TailNumber, aircraft.AltBaro, distance)
}

// PrintAircraft prints a list of Aircraft in a readable manner.
func PrintAircraft(aircraft []Aircraft) {
	if len(aircraft) == 0 {
		fmt.Println("No matching aircraft have been spotted.")
	}

	for _, ac := range aircraft {
		fmt.Println(FormatAircraft(ac))
	}
}
