package jetspotter

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"jetspotter/internal/configuration"
	"jetspotter/internal/metrics"
	"jetspotter/internal/planespotter"
	"jetspotter/internal/weather"

	"github.com/jftuga/geodist"
)

// Vars
var (
	baseURL     = "https://api.adsb.one/v2"
	baseInfoURL = "https://api.adsbdb.com/v0"
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

// getFlightRoute returns the extra information about an aircraft.
func getFlightRoute(callsign string) (route *FlightRoute, err error) {
	endpoint, err := url.JoinPath(baseInfoURL, "callsign", callsign)
	if err != nil {
		return nil, fmt.Errorf("failed to build URL: %w", err)
	}

	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer res.Body.Close()

	// Check HTTP status code
	if res.StatusCode >= 400 {
		if res.StatusCode == 404 {
			return nil, nil
		}

		if res.StatusCode == 429 {
			return nil, fmt.Errorf("API rate limit exceeded: %s", res.Status)
		}

		return nil, fmt.Errorf("API call to %s returned error: %s", endpoint, res.Status)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Use the FlightRouteResponse type from types.go
	var resp FlightRouteResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse API response: %w", err)
	}

	return &resp.Response.FlightRoute, nil
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

// HandleAircraft return a list of aircraft that have been filtered by range, type and altitude.
// Aircraft that have been spotted are removed from the list.
func HandleAircraft(alreadySpottedAircraft *[]Aircraft, config configuration.Config) (aircraft []AircraftOutput, err error) {
	var newlySpottedAircraft []Aircraft

	// Use MaxScanRangeKilometers for scanning (API query)
	allAircraftInRange, err := getAllAircraftInRange(config.Location, config.MaxScanRangeKilometers)
	if err != nil {
		return nil, err
	}

	// Filter the aircraft by the notification range (MaxRangeKilometers)
	var aircraftInNotificationRange []Aircraft
	for _, ac := range allAircraftInRange {
		// Skip aircraft without registration
		if ac.Registration == "" {
			continue
		}

		distance := CalculateDistance(config.Location, geodist.Coord{Lat: ac.Lat, Lon: ac.Lon})
		if distance <= config.MaxRangeKilometers {
			aircraftInNotificationRange = append(aircraftInNotificationRange, ac)
		}
	}

	// For notifications, we need to track what's new and filter by type
	newlySpottedAircraft, *alreadySpottedAircraft = validateAircraft(aircraftInNotificationRange, alreadySpottedAircraft)

	// Only filter by types for notifications, not for the full output
	filteredForNotifications := filterAircraftByTypes(newlySpottedAircraft, config.AircraftTypes)

	// Apply altitude filter if configured
	if config.MaxAltitudeFeet > 0 {
		filteredForNotifications = filterAircraftByAltitude(filteredForNotifications, config.MaxAltitudeFeet)
	}

	// Process newly spotted aircraft for metrics
	newlySpottedAircraftOutput, err := CreateAircraftOutput(newlySpottedAircraft, config, false)
	if err != nil {
		return nil, err
	}
	handleMetrics(newlySpottedAircraftOutput)

	// Generate output for notifications (filtered by AIRCRAFT_TYPES)
	notificationOutputs, err := CreateAircraftOutput(filteredForNotifications, config, true)
	if err != nil {
		return nil, err
	}

	// Update the SpottedAircraft for the API to access - always store ALL aircraft in range
	SpottedAircraft.Lock()
	SpottedAircraft.Aircraft = allAircraftInRange
	SpottedAircraft.Unlock()

	// Return the filtered aircraft for notifications
	return notificationOutputs, nil
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

// filterAircraftByAltitude returns a list of Aircraft that are below the maxAltitudeFeet.
func filterAircraftByAltitude(aircraft []Aircraft, maxAltitudeFeet int) []Aircraft {
	var filteredAircraft []Aircraft

	for _, ac := range aircraft {
		// First convert any 'ground' string indicators to float64(0)
		if ac.AltBaro == "groundft" || ac.AltBaro == "ground" {
			ac.AltBaro = float64(0)
		}

		if ac.AltBaro != nil {
			// Ensure we can safely convert to float64
			var altitude float64
			switch v := ac.AltBaro.(type) {
			case float64:
				altitude = v
			case int:
				altitude = float64(v)
			default:
				// Skip aircraft with unhandled altitude type
				continue
			}

			if int(altitude) <= maxAltitudeFeet {
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
	if aircraft.Callsign == "" || strings.HasPrefix(aircraft.Callsign, " ") {
		aircraft.Callsign = "UNKNOWN"
	}

	aircraft.Callsign = strings.TrimSpace(aircraft.Callsign)

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
// Specify true for extraInfo to include additional information such as flight route, origin, and destination.
func CreateAircraftOutput(aircraft []Aircraft, config configuration.Config, extraInfo bool) (acOutputs []AircraftOutput, err error) {
	var acOutput AircraftOutput
	cloudForecastSucceeded := true

	weather, err := weather.GetCloudForecast(config.Location)
	if err != nil {
		log.Printf("Error getting cloud forecast: %v\n", err)
		cloudForecastSucceeded = false
	}

	for _, ac := range aircraft {
		// Reset aircraft output for new aircraft
		acOutput = AircraftOutput{} // Reset to empty to prevent data leakage between iterations

		ac = validateFields(ac)
		aircraftLocation := geodist.Coord{Lat: ac.Lat, Lon: ac.Lon}
		image := planespotter.GetImageFromAPI(ac.ICAO, ac.Registration)

		acOutput.Altitude = ac.AltBaro.(float64)
		acOutput.Callsign = ac.Callsign
		acOutput.Description = ac.Desc
		acOutput.Distance = CalculateDistance(config.Location, aircraftLocation)
		acOutput.Speed = int(ac.GS)
		acOutput.Registration = ac.Registration
		acOutput.Country = GetCountryFromRegistration(ac.Registration)
		acOutput.Type = ac.PlaneType
		acOutput.ICAO = ac.ICAO
		acOutput.Heading = ac.Track
		acOutput.TrackerURL = fmt.Sprintf("https://globe.airplanes.live/?icao=%v&SiteLat=%f&SiteLon=%f&zoom=11&enableLabels&extendedLabels=1&noIsolation",
			ac.ICAO, config.Location.Lat, config.Location.Lon)
		if cloudForecastSucceeded {
			acOutput.CloudCoverage = getCloudCoverage(*weather, acOutput.Altitude)
		}
		acOutput.BearingFromLocation = CalculateBearing(config.Location, aircraftLocation)
		acOutput.BearingFromAircraft = CalculateBearing(aircraftLocation, config.Location)
		if image != nil {
			acOutput.ImageThumbnailURL = image.ThumbnailLarge.Src
			acOutput.ImageURL = image.Link
			acOutput.Photographer = image.Photographer
		}
		acOutput.Military = isAircraftMilitary(ac)
		// Check if aircraft is on the ground (altitude is 0)
		acOutput.OnGround = acOutput.Altitude == 0
		// If the aircraft is on the ground, it cannot be inbound
		if acOutput.OnGround {
			acOutput.Inbound = false
		} else {
			acOutput.Inbound = IsAircraftInbound(config.Location, ac, 30)
		}

		if extraInfo && ac.Callsign != "UNKNOWN" && len(ac.Callsign) > 3 {
			// Fetch flight route information
			flightRoute, err := getFlightRoute(ac.Callsign)
			if err == nil && flightRoute != nil {
				// Validate that the flight route matches this aircraft
				if isValidFlightRouteForAircraft(flightRoute, ac) {
					acOutput.Airline = flightRoute.Airline
					acOutput.Origin = flightRoute.Origin
					acOutput.Destination = flightRoute.Destination
				} else {
					// Flight route doesn't match this aircraft, log a message for debugging
					log.Printf("Flight route for callsign %s doesn't match aircraft (ICAO: %s, Reg: %s)",
						ac.Callsign, ac.ICAO, ac.Registration)
				}
			} else if err != nil && !strings.Contains(err.Error(), "API rate limit exceeded") {
				// Only log errors that aren't rate limit related
				log.Printf("Error getting flight route information for %s: %v", ac.Callsign, err)
			}
		}

		acOutputs = append(acOutputs, acOutput)
	}
	return acOutputs, nil
}

// isValidFlightRouteForAircraft validates whether a flight route likely matches the given aircraft
// by checking for ICAO/IATA code consistency, registration consistency, or other relevant factors
func isValidFlightRouteForAircraft(route *FlightRoute, aircraft Aircraft) bool {
	// Skip validation if we're missing essential data
	if route == nil {
		return false
	}

	// Basic validation: ensure origin and destination aren't empty
	if route.Origin.Name == "" || route.Destination.Name == "" {
		return false
	}

	// Check if the airline ICAO code prefix matches the aircraft callsign prefix
	// Many airlines use their ICAO code as the callsign prefix
	if route.Airline.ICAO != "" && len(aircraft.Callsign) >= 3 {
		// Many airlines' callsigns start with their ICAO code
		airlinePrefix := strings.ToUpper(aircraft.Callsign[:3])
		if strings.EqualFold(airlinePrefix, route.Airline.ICAO) {
			return true
		}
	}

	// Check if registration matches callsign (some private flights use registration as callsign)
	if aircraft.Registration != "" {
		regWithoutHyphen := strings.ReplaceAll(aircraft.Registration, "-", "")
		if strings.Contains(strings.ToUpper(aircraft.Callsign), regWithoutHyphen) {
			return true
		}
	}

	// Check if callsign is in expected format for scheduled flights
	// Most commercial flights use a format like "ABC123" where ABC is airline code
	if len(aircraft.Callsign) >= 5 && len(route.CallsignICAO) >= 3 &&
		strings.ToUpper(aircraft.Callsign[:3]) == strings.ToUpper(route.CallsignICAO[:3]) {
		return true
	}

	// If all validation checks fail, consider it a mismatch
	return false
}

// IsAircraftInbound checks if the aircraft is inbound to the target location
func IsAircraftInbound(location geodist.Coord, aircraft Aircraft, margin float64) bool {
	// Calculate bearing from aircraft to location (where the aircraft should be pointing if heading to target)
	bearingFromAircraft := CalculateBearing(geodist.Coord{Lat: aircraft.Lat, Lon: aircraft.Lon}, location)

	// Calculate the absolute difference between the ideal bearing and actual aircraft heading
	diff := math.Abs(bearingFromAircraft - aircraft.Track)

	// If the difference is greater than 180 degrees, take the shorter angle
	if diff > 180 {
		diff = 360 - diff
	}

	// If the difference is within the margin, the aircraft is heading toward the target
	return diff <= margin
}

// ConvertAircraftToOutput converts a slice of Aircraft to a slice of AircraftOutput
func ConvertAircraftToOutput(aircraft []Aircraft) []AircraftOutput {
	config, err := configuration.GetConfig()
	if err != nil {
		log.Printf("Error getting config for API: %v", err)
		return []AircraftOutput{}
	}

	// Filter out aircraft without registration
	var filteredAircraft []Aircraft
	for _, ac := range aircraft {
		if ac.Registration != "" {
			filteredAircraft = append(filteredAircraft, ac)
		}
	}

	outputs, err := CreateAircraftOutput(filteredAircraft, config, true)
	if err != nil {
		log.Printf("Error creating aircraft output for API: %v", err)
		return []AircraftOutput{}
	}

	return outputs
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

// GetCountryFromRegistration determines the country of an aircraft based on the registration prefix
func GetCountryFromRegistration(registration string) string {
	if registration == "" {
		return "Unknown"
	}

	// Map common registration prefixes to countries
	registryMap := map[string]string{
		// A few examples from each region
		// North America
		"87-": "United States", // Lockheed C-5M
		"N":   "United States",
		"C-F": "Canada",
		"C-G": "Canada",
		"C-I": "Canada",
		"XA":  "Mexico",
		"XB":  "Mexico",
		"XC":  "Mexico",

		// Europe
		"G-":  "United Kingdom",
		"F-":  "France",
		"D-":  "Germany",
		"I-":  "Italy",
		"EC-": "Spain",
		"CS-": "Portugal",
		"EI-": "Ireland",
		"OE-": "Austria",
		"4L-": "Georgia",
		"TF-": "Iceland",
		"LZ-": "Bulgaria",
		"T7-": "San Marino",
		"HB-": "Switzerland",
		"ER-": "Moldova",
		"9A-": "Croatia",
		"ES-": "Estonia",
		"OO-": "Belgium",
		"FA-": "Belgium", // Belgian F16
		"CT-": "Belgium", // Belgian A400
		"ST-": "Belgium", // Belgian AERMACCHI
		"RN-": "Belgium", // Belgian NH-90
		"YL-": "Latvia",
		"PH-": "Netherlands",
		"L-":  "Netherlands", // Dutch PILATUS
		"SE-": "Sweden",
		"OY-": "Denmark",
		"OH-": "Finland",
		"LN-": "Norway",
		"YR-": "Romania",
		"SP-": "Poland",
		"OK-": "Czech Republic",
		"HA-": "Hungary",
		"YU-": "Serbia",
		"LY-": "Lithuania",
		"UR-": "Ukraine",
		"SX-": "Greece",
		"LX-": "Luxembourg",
		"9H-": "Malta",

		// Asia & Oceania
		"JA":  "Japan",
		"B-":  "China",
		"VT-": "India",
		"HS-": "Thailand",
		"PK-": "Indonesia",
		"9M-": "Malaysia",
		"9V-": "Singapore",
		"VH-": "Australia",
		"ZK-": "New Zealand",

		// South America
		"LV-": "Argentina",
		"PP-": "Brazil",
		"PR-": "Brazil",
		"PT-": "Brazil",
		"PU-": "Brazil",
		"CC-": "Chile",
		"HK-": "Colombia",

		// Middle East & Africa
		"4X-":  "Israel",
		"TC-":  "Turkey",
		"SU-":  "Egypt",
		"ZS-":  "South Africa",
		"ET-":  "Ethiopia",
		"5N-":  "Nigeria",
		"7T-":  "Algeria",
		"TS-":  "Tunisia",
		"CN-":  "Morocco",
		"HZ-":  "Saudi Arabia",
		"A6-":  "United Arab Emirates",
		"A7-":  "Qatar",
		"A9C-": "Bahrain",
		"EP-":  "Iran",
		"YI-":  "Iraq",
		"9K-":  "Kuwait",
		"4K-":  "Azerbaijan",
		"9XR-": "Rwanda",
	}

	// Check for matching prefixes
	for prefix, country := range registryMap {
		if strings.HasPrefix(registration, prefix) {
			return country
		}
	}

	// If no match found, return the first character as a basic hint
	if len(registration) > 0 {
		return "Unknown (" + string(registration[0]) + ")"
	}

	return "Unknown"
}
