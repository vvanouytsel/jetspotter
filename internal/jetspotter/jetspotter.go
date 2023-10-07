package jetspotter

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"net/url"
	"slices"
	"sort"
	"strconv"
	"time"

	"jetspotter/internal/configuration"
	"jetspotter/internal/weather"

	"github.com/PuerkitoBio/goquery"
	"github.com/jftuga/geodist"
)

// Vars
var (
	baseURL = "https://api.adsb.one/v2/point"
	// baseURL = "https://api.airplanes.live/v2/point"
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

	if aircraft.AltBaro == "groundft" || aircraft.AltBaro == "ground" || aircraft.AltBaro == nil {
		aircraft.AltBaro = float64(0)
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

	for _, ac := range aircraft {
		ac = validateFields(ac)
		aircraftLocation := geodist.Coord{Lat: ac.Lat, Lon: ac.Lon}

		acOutput.Altitude = ac.AltBaro.(float64)
		acOutput.Callsign = ac.Callsign
		acOutput.Description = ac.Desc
		acOutput.Distance = CalculateDistance(config.Location, aircraftLocation)
		acOutput.Speed = int(ac.GS)
		acOutput.TailNumber = ac.TailNumber
		acOutput.Type = ac.PlaneType
		acOutput.ICAO = ac.ICAO
		acOutput.Heading = ac.Track
		acOutput.TrackerURL = fmt.Sprintf("https://globe.adsbexchange.com/?icao=%s", ac.ICAO)
		acOutput.CloudCoverage = getCloudCoverage(*weather, acOutput.Altitude)
		acOutput.BearingFromLocation = CalculateBearing(config.Location, aircraftLocation)
		acOutput.BearingFromAircraft = CalculateBearing(aircraftLocation, config.Location)
		acOutput.PlaneSpotterURL = getImageURL(fmt.Sprintf("https://www.planespotting.be/index.php?page=aircraft&registration=%s", ac.TailNumber))

		acOutputs = append(acOutputs, acOutput)
	}
	return acOutputs, nil
}

// SortByDistance sorts a slice of aircraft to show the closest aircraft first
func SortByDistance(aircraft []AircraftOutput) []AircraftOutput {
	sort.Slice(aircraft, func(i, j int) bool {
		return aircraft[i].Distance < aircraft[j].Distance
	})

	return aircraft
}

// CalculateBearing returns the bearing from the source to the target
func CalculateBearing(source geodist.Coord, target geodist.Coord) float64 {
	y := math.Sin(toRadians(target.Lon-source.Lon)) * math.Cos(toRadians(target.Lat))
	x := math.Cos(toRadians(source.Lat))*math.Sin(toRadians(target.Lat)) - math.Sin(toRadians(source.Lat))*math.Cos(toRadians(target.Lat))*math.Cos(toRadians(target.Lon-source.Lon))

	bearing := math.Atan2(y, x)
	bearing = (toDegrees(bearing) + 360)

	if bearing >= 360 {
		bearing -= 360
	}

	return bearing
}

func toRadians(degrees float64) float64 {
	return degrees * (math.Pi / 180)
}

func toDegrees(rad float64) float64 {
	return rad * (180 / math.Pi)
}

// getImageURL queries the planespotter website to fetch an image link
func getImageURL(URL string) (imageURL string) {
	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		return ""
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return ""
	}

	if res.StatusCode != 200 {
		fmt.Printf("Received status code %d for URL %s\n", res.StatusCode, URL)
		return ""
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	doc.Find("meta[property='og:image']").Each(func(index int, element *goquery.Selection) {
		if content, exists := element.Attr("content"); exists {
			imageURL = content
		}
	})

	return imageURL
}
