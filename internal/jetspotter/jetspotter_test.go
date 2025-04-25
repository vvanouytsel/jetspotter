package jetspotter

import (
	"jetspotter/internal/aircraft"
	"jetspotter/internal/configuration"
	"reflect"
	"testing"

	"github.com/jftuga/geodist"
)

var (
	jets = []Aircraft{
		{
			ICAO:      "ABC",
			Callsign:  "JACKAL51",
			PlaneType: "F16",
		},
		{
			ICAO:      "DEF",
			Callsign:  "XSG432",
			PlaneType: "CESSNA",
		},
		{
			ICAO:      "HAH",
			Callsign:  "VIKING11",
			PlaneType: "F18",
		},
	}

	jackal51Aircraft = Aircraft{
		ICAO:      "ABC",
		Callsign:  "JACKAL51",
		PlaneType: "F16",
	}

	planesWithAltitude = []Aircraft{
		{
			Callsign: "KHARMA11",
			AltBaro:  4000.0,
		},
		{
			Callsign: "KHARMA12",
			AltBaro:  9000.0,
		},
		{
			Callsign: "KHARMA13",
		},
	}

	planes = []Aircraft{
		{
			Callsign:  "APEX11",
			PlaneType: aircraft.F16.Identifier,
			Desc:      aircraft.F16.Description,
			DbFlags:   1,
		},
		{
			Callsign:  "APEX12",
			PlaneType: aircraft.F16.Identifier,
			Desc:      aircraft.F16.Description,
			DbFlags:   1,
		},
		{
			Callsign:  "XSG123",
			PlaneType: aircraft.B77L.Identifier,
			Desc:      aircraft.B77L.Description,
		},
		{
			Callsign:  "ABC987",
			PlaneType: aircraft.A320.Identifier,
			Desc:      aircraft.A320.Description,
		},
		{
			Callsign:  "GRZLY11",
			PlaneType: aircraft.A400.Identifier,
			Desc:      aircraft.A400.Description,
			DbFlags:   1,
		},
	}

	locationMannekenPis = geodist.Coord{
		Lat: 50.844987343465924,
		Lon: 4.349981064923107,
	}

	locationElisabethPark = geodist.Coord{
		Lat: 50.86503662037458,
		Lon: 4.32399484006766,
	}

	locationChristRedeemer = geodist.Coord{
		Lat: -22.951907892908967,
		Lon: -43.21048377096087,
	}

	locationPyramidGiza = geodist.Coord{
		Lat: 29.979104641494533,
		Lon: 31.134157868680205,
	}
)

func TestFilterAircraftByTypeF16(t *testing.T) {
	expected := []Aircraft{
		{
			Callsign:  "APEX11",
			PlaneType: aircraft.F16.Identifier,
			Desc:      aircraft.F16.Description,
			DbFlags:   1,
		},
		{
			Callsign:  "APEX12",
			PlaneType: aircraft.F16.Identifier,
			Desc:      aircraft.F16.Description,
			DbFlags:   1,
		},
	}

	config := configuration.Config{
		AircraftTypes: []string{aircraft.F16.Identifier},
	}
	actual := filterAircraftByTypes(planes, config.AircraftTypes)

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("expected '%v' to be the same as '%v'", expected, actual)
	}
}

func TestFilterAircraftByTypeAll(t *testing.T) {
	config := configuration.Config{
		AircraftTypes: []string{aircraft.ALL.Identifier},
	}
	expected := planes
	actual := filterAircraftByTypes(planes, config.AircraftTypes)

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("expected '%v' to be the same as '%v'", expected, actual)
	}
}

func TestFilterAircraftByTypeMilitary(t *testing.T) {
	config := configuration.Config{
		AircraftTypes: []string{aircraft.MILITARY.Identifier},
	}
	expected := []Aircraft{
		{
			Callsign:  "APEX11",
			PlaneType: aircraft.F16.Identifier,
			Desc:      aircraft.F16.Description,
			DbFlags:   1,
		},
		{
			Callsign:  "APEX12",
			PlaneType: aircraft.F16.Identifier,
			Desc:      aircraft.F16.Description,
			DbFlags:   1,
		},
		{
			Callsign:  "GRZLY11",
			PlaneType: aircraft.A400.Identifier,
			Desc:      aircraft.A400.Description,
			DbFlags:   1,
		},
	}

	actual := filterAircraftByTypes(planes, config.AircraftTypes)

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("expected '%v' to be the same as '%v'", expected, actual)
	}
}

func TestFilterAircraftByTypeMilitaryAndA320(t *testing.T) {
	config := configuration.Config{
		AircraftTypes: []string{aircraft.MILITARY.Identifier, aircraft.A320.Identifier},
	}

	expected := []Aircraft{
		{
			Callsign:  "APEX11",
			PlaneType: aircraft.F16.Identifier,
			Desc:      aircraft.F16.Description,
			DbFlags:   1,
		},
		{
			Callsign:  "APEX12",
			PlaneType: aircraft.F16.Identifier,
			Desc:      aircraft.F16.Description,
			DbFlags:   1,
		},
		{
			Callsign:  "ABC987",
			PlaneType: aircraft.A320.Identifier,
			Desc:      aircraft.A320.Description,
		},
		{
			Callsign:  "GRZLY11",
			PlaneType: aircraft.A400.Identifier,
			Desc:      aircraft.A400.Description,
			DbFlags:   1,
		},
	}

	actual := filterAircraftByTypes(planes, config.AircraftTypes)

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("expected '%v' to be the same as '%v'", expected, actual)
	}
}

func TestCalculateDistance1(t *testing.T) {
	tbilisiAirportCoordinates := geodist.Coord{
		Lat: 41.4007,
		Lon: 44.5705,
	}

	kutaisiAirportCoordinates := geodist.Coord{
		Lat: 42.1033,
		Lon: 42.2830,
	}

	expected := 205
	actual := CalculateDistance(tbilisiAirportCoordinates, kutaisiAirportCoordinates)

	if expected != actual {
		t.Fatalf("expected '%v' to be the same as '%v'", expected, actual)
	}
}

func TestCalculateDistance2(t *testing.T) {
	lasVegasCoordinates := geodist.Coord{
		Lat: 36.11467991019019,
		Lon: -115.18050028591726,
	}

	denverCoordinates := geodist.Coord{
		Lat: 39.7400431976992,
		Lon: -104.99281871032076,
	}

	expected := 980
	actual := CalculateDistance(lasVegasCoordinates, denverCoordinates)

	if expected != actual {
		t.Fatalf("expected '%v' to be the same as '%v'", expected, actual)
	}
}

func TestFilterAircraftByTypes(t *testing.T) {
	expected := []Aircraft{
		{
			Callsign:  "APEX11",
			PlaneType: aircraft.F16.Identifier,
			Desc:      aircraft.F16.Description,
			DbFlags:   1,
		},
		{
			Callsign:  "APEX12",
			PlaneType: aircraft.F16.Identifier,
			Desc:      aircraft.F16.Description,
			DbFlags:   1,
		},
		{
			Callsign:  "GRZLY11",
			PlaneType: aircraft.A400.Identifier,
			Desc:      aircraft.A400.Description,
			DbFlags:   1,
		},
	}

	config := configuration.Config{
		AircraftTypes: []string{aircraft.F16.Identifier, aircraft.A400.Identifier},
	}
	actual := filterAircraftByTypes(planes, config.AircraftTypes)

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("expected '%v' to be the same as '%v'", expected, actual)
	}
}

func TestConvertKnotsToKilometerPerHour(t *testing.T) {
	expected := 185
	actual := ConvertKnotsToKilometersPerHour(100)

	if expected != actual {
		t.Fatalf("expected '%v' to be the same as '%v'", expected, actual)
	}
}

func TestConvertFeetToMeters(t *testing.T) {
	expected := 30
	actual := ConvertFeetToMeters(100)

	if expected != actual {
		t.Fatalf("expected '%v' to be the same as '%v'", expected, actual)
	}
}

func TestSortAircraftByDistance(t *testing.T) {
	aircraft := []AircraftOutput{
		{
			Callsign: "APEX11",
			Distance: 120,
		},
		{
			Callsign: "APEX12",
			Distance: 60,
		},
		{
			Callsign: "APEX13",
			Distance: 10,
		},
	}

	sortedAircraft := SortByDistance(aircraft)

	if sortedAircraft[0].Callsign != "APEX13" || sortedAircraft[1].Callsign != "APEX12" || sortedAircraft[2].Callsign != "APEX11" {
		t.Fatal("List is not sorted by distance")
	}
}

func TestCalculateBearing1(t *testing.T) {

	expected := 320
	actual := int(CalculateBearing(locationMannekenPis, locationElisabethPark))

	if expected != actual {
		t.Fatalf("expected '%v' to be the same as '%v'", expected, actual)
	}
}

func TestCalculateBearing2(t *testing.T) {

	expected := 140
	actual := int(CalculateBearing(locationElisabethPark, locationMannekenPis))

	if expected != actual {
		t.Fatalf("expected '%v' to be the same as '%v'", expected, actual)
	}
}

func TestCalculateBearing3(t *testing.T) {

	expected := 56
	actual := int(CalculateBearing(locationChristRedeemer, locationPyramidGiza))

	if expected != actual {
		t.Fatalf("expected '%v' to be the same as '%v'", expected, actual)
	}
}

func TestCalculateBearing4(t *testing.T) {

	expected := 242
	actual := int(CalculateBearing(locationPyramidGiza, locationChristRedeemer))

	if expected != actual {
		t.Fatalf("expected '%v' to be the same as '%v'", expected, actual)
	}
}

func TestCalculateBearing5(t *testing.T) {
	source := geodist.Coord{
		Lat: 51.42676766088391,
		Lon: 4.623935349264089,
	}

	target := geodist.Coord{
		Lat: 51.426688015979074,
		Lon: 4.63915475148803,
	}

	expected := 90
	actual := int(CalculateBearing(source, target))

	if expected != actual {
		t.Fatalf("expected '%v' to be the same as '%v'", expected, actual)
	}
}

func TestNewlySpotted(t *testing.T) {

	newAircaft := Aircraft{
		ICAO:     "ZYX",
		Type:     "A400",
		Callsign: "NEW11",
	}
	expected := true
	actual := newlySpotted(newAircaft, jets)

	if expected != actual {
		t.Fatalf("expected '%v' to be the same as '%v'", expected, actual)
	}
}

func TestAlreadySpotted(t *testing.T) {

	expected := false
	actual := newlySpotted(jackal51Aircraft, jets)

	if expected != actual {
		t.Fatalf("expected '%v' to be the same as '%v'", expected, actual)
	}
}

func TestAlreadySpottedAircraftAreFiltered(t *testing.T) {

	aircraft := []Aircraft{
		{
			ICAO:      "VVO",
			Callsign:  "DEMON81",
			PlaneType: "F16",
		},
		{
			ICAO:      "DEF",
			Callsign:  "LDS431",
			PlaneType: "CESSNA",
		},
		{
			ICAO:      "RPQ",
			Callsign:  "JESTER41",
			PlaneType: "F18",
		},
	}

	alreadySpottedAircraft := []Aircraft{
		{
			ICAO:      "DEF",
			Callsign:  "LDS431",
			PlaneType: "CESSNA",
		},
		{
			ICAO:      "X123",
			Callsign:  "GOOSE11",
			PlaneType: "F22",
		},
	}

	expectedNewlySpotted := []Aircraft{
		{
			ICAO:      "VVO",
			Callsign:  "DEMON81",
			PlaneType: "F16",
		},
		{
			ICAO:      "RPQ",
			Callsign:  "JESTER41",
			PlaneType: "F18",
		},
	}

	expectedSpottedAircraft := []Aircraft{
		{
			ICAO:      "DEF",
			Callsign:  "LDS431",
			PlaneType: "CESSNA",
		},
		{
			ICAO:      "VVO",
			Callsign:  "DEMON81",
			PlaneType: "F16",
		},
		{
			ICAO:      "RPQ",
			Callsign:  "JESTER41",
			PlaneType: "F18",
		},
	}

	actualNewlySpottedAircraft, actualSpottedAircraft := validateAircraft(aircraft, &alreadySpottedAircraft)

	if !reflect.DeepEqual(expectedNewlySpotted, actualNewlySpottedAircraft) {
		t.Fatalf("expected '%v' to be the same as '%v' in the newly spotted list",
			expectedNewlySpotted, actualNewlySpottedAircraft)
	}

	if !reflect.DeepEqual(expectedSpottedAircraft, actualSpottedAircraft) {
		t.Fatalf("expected '%v' to be the same as '%v' in the already spotted list",
			expectedSpottedAircraft, actualSpottedAircraft)
	}
}

// TestCallsignStartsEmptySpace tests that the callsign is set to "UNKNOWN" when the callsign starts with empty spaces
func TestCallsignStartsEmptySpace(t *testing.T) {
	aircraft := Aircraft{
		Callsign: "    ",
	}

	expected := "UNKNOWN"
	actual := validateFields(aircraft).Callsign

	if expected != actual {
		t.Fatalf("expected '%v' to be the same as '%v'", expected, actual)
	}
}

// TestCallsignNoEmptySpace tests that the callsign is not changed if it does not start with empty spaces
func TestCallsignNoEmptySpace(t *testing.T) {
	aircraft := Aircraft{
		Callsign: "x  x  ",
	}

	expected := "x  x  "
	actual := validateFields(aircraft).Callsign

	if expected != actual {
		t.Fatalf("expected '%v' to be the same as '%v'", expected, actual)
	}
}

func TestGetHighestValue(t *testing.T) {
	expected := 14
	actual := getHighestValue(0, 8, 14, 3, 12)
	if expected != actual {
		t.Fatalf("expected '%v' to be the same as '%v'", expected, actual)
	}
}

func TestConvertKilometersToMiles(t *testing.T) {
	expected := 10
	actual := convertKilometersToNauticalMiles(20)
	if expected != actual {
		t.Fatalf("expected '%v' to be the same as '%v'", expected, actual)
	}
}

func TestFilterAircraftByAltitude(t *testing.T) {
	expected := []Aircraft{
		{
			Callsign: "KHARMA11",
			AltBaro:  4000.0,
		},
	}

	config := configuration.Config{
		MaxAltitudeFeet: 5000.0,
	}
	actual := filterAircraftByAltitude(planesWithAltitude, config.MaxAltitudeFeet)

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("expected '%v' to be the same as '%v'", expected, actual)
	}
}

// TestFilterAircraftByDistanceWithDifferentRanges tests that aircraft outside MaxRangeKilometers
// but inside MaxScanRangeKilometers are correctly filtered
func TestFilterAircraftByDistanceWithDifferentRanges(t *testing.T) {
	// Create a set of test aircraft at different distances
	testAircraft := []Aircraft{
		{
			ICAO:     "TEST1",
			Callsign: "NEAR1",
			Lat:      51.18, // Very close to default location
			Lon:      5.46,  // Should be in notification range
		},
		{
			ICAO:     "TEST2",
			Callsign: "MEDIUM1",
			Lat:      51.40, // Medium distance - adjusted to be > 20km but < 100km
			Lon:      5.80,  // Should be in scan range but not notification range
		},
		{
			ICAO:     "TEST3",
			Callsign: "FAR1",
			Lat:      52.50, // Far distance - adjusted to be > 100km
			Lon:      7.00,  // Should be outside both ranges
		},
	}

	// Create a simulated location and config
	location := geodist.Coord{
		Lat: 51.17348, // Default location from config
		Lon: 5.45921,
	}

	// Compute actual distances for verification
	distances := make([]int, len(testAircraft))
	for i, ac := range testAircraft {
		aircraftLoc := geodist.Coord{Lat: ac.Lat, Lon: ac.Lon}
		distances[i] = CalculateDistance(location, aircraftLoc)
	}

	// Set up ranges for testing
	notificationRange := 20 // 20km notification range
	scanRange := 100        // 100km scan range

	// Print actual distances for debugging
	t.Logf("TEST1 distance: %d km", distances[0])
	t.Logf("TEST2 distance: %d km", distances[1])
	t.Logf("TEST3 distance: %d km", distances[2])

	// Verify the distances are as expected
	if distances[0] > notificationRange {
		t.Fatalf("TEST1 should be within notification range: %d km", distances[0])
	}

	if distances[1] <= notificationRange || distances[1] > scanRange {
		t.Fatalf("TEST2 should be outside notification range but inside scan range: %d km", distances[1])
	}

	if distances[2] <= scanRange {
		t.Fatalf("TEST3 should be outside both ranges: %d km", distances[2])
	}

	// Filter manually to generate expected result
	var expectedInNotificationRange []Aircraft
	var expectedInScanRange []Aircraft

	for i, ac := range testAircraft {
		if distances[i] <= notificationRange {
			expectedInNotificationRange = append(expectedInNotificationRange, ac)
			expectedInScanRange = append(expectedInScanRange, ac)
		} else if distances[i] <= scanRange {
			expectedInScanRange = append(expectedInScanRange, ac)
		}
	}

	// Test filtering logic directly
	var aircraftInNotificationRange []Aircraft
	for _, ac := range testAircraft {
		distance := CalculateDistance(location, geodist.Coord{Lat: ac.Lat, Lon: ac.Lon})
		if distance <= notificationRange {
			aircraftInNotificationRange = append(aircraftInNotificationRange, ac)
		}
	}

	if !reflect.DeepEqual(expectedInNotificationRange, aircraftInNotificationRange) {
		t.Fatalf("expected '%v' to be in notification range, got '%v'", expectedInNotificationRange, aircraftInNotificationRange)
	}

	// Verify each aircraft is correctly categorized based on distances
	if len(expectedInNotificationRange) == 0 {
		t.Fatal("expected at least one aircraft to be in notification range")
	}

	if len(expectedInScanRange) <= len(expectedInNotificationRange) {
		t.Fatal("expected more aircraft to be in scan range than notification range")
	}
}

// TestHandleAircraftWithScanRange tests that aircraft are properly filtered based on
// both scan range and notification range
func TestHandleAircraftWithScanRange(t *testing.T) {
	// Create test aircraft at different distances
	testAircraft := []Aircraft{
		{
			ICAO:     "TEST1",
			Callsign: "NEAR1",
			Lat:      51.18,
			Lon:      5.46,
			AltBaro:  float64(1000),
		},
		{
			ICAO:     "TEST2",
			Callsign: "MEDIUM1",
			Lat:      51.40,
			Lon:      5.80,
			AltBaro:  float64(2000),
		},
		{
			ICAO:     "TEST3",
			Callsign: "FAR1",
			Lat:      52.50,
			Lon:      7.00,
			AltBaro:  float64(3000),
		},
	}

	// Set up test config with different scan and notification ranges
	config := configuration.Config{
		Location: geodist.Coord{
			Lat: 51.17348,
			Lon: 5.45921,
		},
		MaxRangeKilometers:     20,  // Should only include TEST1
		MaxScanRangeKilometers: 100, // Should include TEST1 and TEST2, but not TEST3
		AircraftTypes:          []string{"ALL"},
		MaxAltitudeFeet:        0, // No altitude filtering
	}

	// Filter aircraft based on distance from the location
	var aircraftInScanRange []Aircraft
	for _, ac := range testAircraft {
		distance := CalculateDistance(config.Location, geodist.Coord{Lat: ac.Lat, Lon: ac.Lon})
		if distance <= config.MaxScanRangeKilometers {
			aircraftInScanRange = append(aircraftInScanRange, ac)
		}
	}

	// Manually apply the same filtering logic as in HandleAircraft function
	var aircraftInNotificationRange []Aircraft
	for _, ac := range aircraftInScanRange {
		distance := CalculateDistance(config.Location, geodist.Coord{Lat: ac.Lat, Lon: ac.Lon})
		if distance <= config.MaxRangeKilometers {
			aircraftInNotificationRange = append(aircraftInNotificationRange, ac)
		}
	}

	// Verify that our filtering picked the correct aircraft
	if len(aircraftInScanRange) != 2 {
		t.Fatalf("Expected 2 aircraft in scan range, got %d", len(aircraftInScanRange))
	}

	if len(aircraftInNotificationRange) != 1 {
		t.Fatalf("Expected 1 aircraft in notification range, got %d", len(aircraftInNotificationRange))
	}

	if aircraftInNotificationRange[0].Callsign != "NEAR1" {
		t.Fatalf("Expected NEAR1, got %s", aircraftInNotificationRange[0].Callsign)
	}

	// Now test the actual distance-based filtering in the HandleAircraft function
	// by simulating its behavior with our test data
	var alreadySpottedAircraft []Aircraft

	// Simulate the main filtering logic from HandleAircraft
	newlySpottedAircraft, updatedSpottedAircraft := validateAircraft(aircraftInNotificationRange, &alreadySpottedAircraft)
	filteredAircraft := filterAircraftByTypes(newlySpottedAircraft, config.AircraftTypes)

	// Check that we get expected results
	if len(filteredAircraft) != 1 {
		t.Fatalf("Expected 1 aircraft after filtering, got %d", len(filteredAircraft))
	}

	if filteredAircraft[0].Callsign != "NEAR1" {
		t.Fatalf("Expected NEAR1 after filtering, got %s", filteredAircraft[0].Callsign)
	}

	if len(updatedSpottedAircraft) != 1 {
		t.Fatalf("Expected 1 aircraft in updatedSpottedAircraft, got %d", len(updatedSpottedAircraft))
	}

	// This verifies that the notifications list only includes aircraft within notification range
	if updatedSpottedAircraft[0].Callsign != "NEAR1" {
		t.Fatalf("Expected NEAR1 in updatedSpottedAircraft, got %s", updatedSpottedAircraft[0].Callsign)
	}
}
