package jetspotter

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/url"
	"slices"
	"sort"
	"strconv"
	"time"

	"jetspotter/internal/configuration"
	"jetspotter/internal/metrics"
	"jetspotter/internal/planespotter"
	"jetspotter/internal/weather"

	"github.com/jftuga/geodist"
)

// Vars
var (
	baseURL = "https://api.adsb.one/v2"
	// baseURL = "https://api.airplanes.live/v2"
)

// CalculateDistance returns the rounded distance between two coordinates in kilometers
func CalculateDistance(source geodist.Coord, destination geodist.Coord) int {
	_, kilometers := geodist.HaversineDistance(source, destination)
	return int(kilometers)
}

// convertKilometersToNauticalMiles converts kilometers into miles. The miles are rounded.
func convertKilometersToNauticalMiles(kilometers float64) int {
	return int(kilometers / 1.852)
}

// getMilitaryAircraftInRange gets all the military aircraft on the map, loops over each aircraft and returns
// only the aircraft that are within the specified maxRangeKilometers.
func getMilitaryAircraftInRange(location geodist.Coord, maxRangeKilometers int) (aircraft []Aircraft, err error) {
	var flightData FlightData
	endpoint, err := url.JoinPath(baseURL, "mil")
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

	for _, ac := range flightData.AC {
		distance := CalculateDistance(location, geodist.Coord{Lat: ac.Lat, Lon: ac.Lon})
		if distance <= maxRangeKilometers {
			aircraft = append(aircraft, ac)
		}
	}
	return aircraft, nil
}

// getAllAircraftInRange returns all aircraft within maxRange kilometers of the location.
func getAllAircraftInRange(location geodist.Coord, maxRangeKilometers int) (aircraft []Aircraft, err error) {
	var flightData FlightData
	miles := convertKilometersToNauticalMiles(float64(maxRangeKilometers))
	endpoint, err := url.JoinPath(baseURL, "point",
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

	return flightData.AC, nil
}

// newlySpotted returns true if the aircraft has not been spotted during the last interval.
func newlySpotted(aircraft Aircraft, spottedAircraft []Aircraft) bool {
	return !containsAircraft(aircraft, spottedAircraft)
}

// containsAircraft checks if the aircraft exists in the list of aircraft.
func containsAircraft(aircraft Aircraft, aircraftList []Aircraft) bool {
	for _, ac := range aircraftList {
		if ac.ICAO == aircraft.ICAO {
			return true
		}
	}
	return false
}

// updateSpottedAircraft removed the previously spotted aircraft that are no longer in range.
func updateSpottedAircraft(alreadySpottedAircraft, filteredAircraft []Aircraft) (aircraft []Aircraft) {
	for _, ac := range alreadySpottedAircraft {
		if containsAircraft(ac, filteredAircraft) {
			aircraft = append(aircraft, ac)
		}
	}

	return aircraft
}

// validateAircraft returns a list of aircraft that have not yet been spotted and
// a list of aircraft that are already spotted, aircraft that were previously spotted but haven't been spotted
// in the last attempt are removed from the already spotted list.
// In practice this means that if an aircraft leaves the spotting range, it is removed from the already spotted list
// and thus the next time they appear in range, a notification will be sent for that aircraft.
func validateAircraft(allFilteredAircraft []Aircraft, alreadySpottedAircraft *[]Aircraft) (newlySpottedAircraft, updatedSpottedAircraft []Aircraft) {
	for _, ac := range allFilteredAircraft {
		if newlySpotted(ac, *alreadySpottedAircraft) {
			newlySpottedAircraft = append(newlySpottedAircraft, ac)
			*alreadySpottedAircraft = append(*alreadySpottedAircraft, ac)
		}
	}

	*alreadySpottedAircraft = updateSpottedAircraft(*alreadySpottedAircraft, allFilteredAircraft)
	return newlySpottedAircraft, *alreadySpottedAircraft
}

// HandleAircraft return a list of aircraft that have been filtered by range and type.
// Aircraft that have been spotted are removed from the list.
func HandleAircraft(alreadySpottedAircraft *[]Aircraft, config configuration.Config) (aircraft []AircraftOutput, err error) {
	var newlySpottedAircraft []Aircraft

	allAircraftInRange, err := getAllAircraftInRange(config.Location, config.MaxRangeKilometers)
	if err != nil {
		return nil, err
	}

	newlySpottedAircraft, *alreadySpottedAircraft = validateAircraft(allAircraftInRange, alreadySpottedAircraft)
	filteredAircraftInRange := filterAircraftByTypes(newlySpottedAircraft, config.AircraftTypes)
	newlySpottedAircraftOutput, err := CreateAircraftOutput(newlySpottedAircraft, config)
	if err != nil {
		return nil, err
	}
	handleMetrics(newlySpottedAircraftOutput)

	acOutputs, err := CreateAircraftOutput(filteredAircraftInRange, config)
	if err != nil {
		return nil, err
	}

	if slices.Contains(config.AircraftTypes, "ALL") {
		return acOutputs, nil
	}
	return acOutputs, nil
}

func handleMetrics(aircraft []AircraftOutput) {
	for _, ac := range aircraft {
		metrics.IncrementMetrics(ac.Type, ac.Description, strconv.FormatBool(ac.Military), ac.Altitude)
	}
}

func isAircraftMilitary(aircraft Aircraft) bool {
	return aircraft.DbFlags == 1
}

func isAircraftDesired(aircraft Aircraft, aircraftType string) bool {
	if aircraftType == "MILITARY" && aircraft.DbFlags == 1 {
		return true
	}

	if aircraft.PlaneType == aircraftType || aircraftType == "ALL" {
		return true
	}

	return false
}

// filterAircraftByTypes returns a list of Aircraft that match the aircraftTypes.
func filterAircraftByTypes(aircraft []Aircraft, types []string) []Aircraft {
	var filteredAircraft []Aircraft

	for _, ac := range aircraft {
		for _, aircraftType := range types {
			if isAircraftDesired(ac, aircraftType) {
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
func getCloudCoverage(weather weather.Data, altitudeInFeet float64) (cloudCoveragePercentage int) {

	altitudeInMeters := ConvertFeetToMeters(altitudeInFeet)
	hourUTC := (time.Now().Hour())

	switch {
	case altitudeInMeters < 3000:
		return weather.Hourly.CloudcoverLow[hourUTC]
	case altitudeInMeters >= 3000 && altitudeInMeters < 8000:
		return getHighestValue(weather.Hourly.CloudcoverLow[hourUTC], weather.Hourly.CloudcoverMid[hourUTC])
	default:
		return getHighestValue(weather.Hourly.CloudcoverLow[hourUTC],
			weather.Hourly.CloudcoverMid[hourUTC],
			weather.Hourly.CloudcoverHigh[hourUTC])
	}
}

func getHighestValue(numbers ...int) (highest int) {
	highest = 0
	for _, v := range numbers {
		if v > highest {
			highest = v
		}
	}
	return highest
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
		image := planespotter.GetImageFromAPI(ac.ICAO, ac.Registration)

		acOutput.Altitude = ac.AltBaro.(float64)
		acOutput.Callsign = ac.Callsign
		acOutput.Description = ac.Desc
		acOutput.Distance = CalculateDistance(config.Location, aircraftLocation)
		acOutput.Speed = int(ac.GS)
		acOutput.Registration = ac.Registration
		acOutput.Type = ac.PlaneType
		acOutput.ICAO = ac.ICAO
		acOutput.Heading = ac.Track
		acOutput.TrackerURL = fmt.Sprintf("https://globe.adsbexchange.com/?icao=%v&SiteLat=%f&SiteLon=%f&zoom=11&enableLabels&extendedLabels=1&noIsolation",
			ac.ICAO, config.Location.Lat, config.Location.Lon)
		acOutput.CloudCoverage = getCloudCoverage(*weather, acOutput.Altitude)
		acOutput.BearingFromLocation = CalculateBearing(config.Location, aircraftLocation)
		acOutput.BearingFromAircraft = CalculateBearing(aircraftLocation, config.Location)
		acOutput.ImageThumbnailURL = image.ThumbnailLarge.Src
		acOutput.ImageURL = image.Link
		acOutput.Military = isAircraftMilitary(ac)
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
