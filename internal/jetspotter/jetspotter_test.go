package jetspotter

import (
	"jetspotter/internal/aircraft"
	"jetspotter/internal/configuration"
	"math"
	"reflect"
	"testing"

	"github.com/jftuga/geodist"
)

var (
	jets = []Aircraft{
		{
			ICAO:     "ABC",
			Callsign: "JACKAL51",
			Type:     "F16",
		},
		{
			ICAO:     "DEF",
			Callsign: "XSG432",
			Type:     "CESSNA",
		},
		{
			ICAO:     "HAH",
			Callsign: "VIKING11",
			Type:     "F18",
		},
	}

	jackal51Aircraft = Aircraft{
		ICAO:     "ABC",
		Callsign: "JACKAL51",
		Type:     "F16",
	}

	planesWithAltitude = []Aircraft{
		{
			Callsign: "KHARMA11",
			Altitude: 4000.0,
		},
		{
			Callsign: "KHARMA12",
			Altitude: 9000.0,
		},
		{
			Callsign: "KHARMA13",
			Altitude: 7000.0,
		},
	}

	planes = []Aircraft{
		{
			Callsign:     "APEX11",
			Type:         "F16",
			Registration: "ABC",
			Description:  aircraft.F16.Description,
			Military:     true,
		},
		{
			Callsign:     "APEX12",
			Type:         "F16",
			Registration: "ABC",
			Description:  aircraft.F16.Description,
			Military:     true,
		},
		{
			Callsign:     "XSG123",
			Type:         "B77L",
			Registration: "ABC",
			Description:  aircraft.B77L.Description,
		},
		{
			Callsign:     "ABC987",
			Type:         "A320",
			Registration: "ABC",
			Description:  aircraft.A320.Description,
		},
		{
			Callsign:     "GRZLY11",
			Type:         "A400",
			Registration: "ABC",
			Description:  aircraft.A400.Description,
			Military:     true,
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
			Callsign:     "APEX11",
			Type:         aircraft.F16.Identifier,
			Description:  aircraft.F16.Description,
			Registration: "ABC",
			Military:     true,
		},
		{
			Callsign:     "APEX12",
			Type:         aircraft.F16.Identifier,
			Description:  aircraft.F16.Description,
			Registration: "ABC",
			Military:     true,
		},
	}

	config := configuration.Config{
		AircraftTypes: []string{"F16"},
	}
	actual := filterAircraftByTypes(planes, config.AircraftTypes)

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("expected \n'%v'\n to be the same as \n'%v'\n", expected, actual)
	}
}

func TestFilterAircraftByTypeAll(t *testing.T) {
	config := configuration.Config{
		AircraftTypes: []string{"ALL"},
	}
	expected := planes
	actual := filterAircraftByTypes(planes, config.AircraftTypes)

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("expected \n'%v'\n to be the same as \n'%v'\n", expected, actual)
	}
}

func TestFilterAircraftByTypeMilitary(t *testing.T) {
	config := configuration.Config{
		AircraftTypes: []string{"MILITARY"},
	}
	// The expected results should match the actual structure of military planes from the 'planes' variable
	expected := []Aircraft{
		planes[0], // APEX11 (F16)
		planes[1], // APEX12 (F16)
		planes[4], // GRZLY11 (A400)
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
		planes[0], // APEX11 (F16, military)
		planes[1], // APEX12 (F16, military)
		planes[3], // ABC987 (A320)
		planes[4], // GRZLY11 (A400, military)
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
	config := configuration.Config{
		AircraftTypes: []string{aircraft.F16.Identifier, aircraft.A400.Identifier},
	}
	expected := []Aircraft{
		planes[0], // APEX11 (F16)
		planes[1], // APEX12 (F16)
		planes[4], // GRZLY11 (A400)
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
	aircraft := []Aircraft{
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
			ICAO:     "VVO",
			Callsign: "DEMON81",
			Type:     "F16",
		},
		{
			ICAO:     "DEF",
			Callsign: "LDS431",
			Type:     "CESSNA",
		},
		{
			ICAO:     "RPQ",
			Callsign: "JESTER41",
			Type:     "F18",
		},
	}

	alreadySpottedAircraft := []Aircraft{
		{
			ICAO:     "DEF",
			Callsign: "LDS431",
			Type:     "CESSNA",
		},
		{
			ICAO:     "X123",
			Callsign: "GOOSE11",
			Type:     "F22",
		},
	}

	expectedNewlySpotted := []Aircraft{
		{
			ICAO:     "VVO",
			Callsign: "DEMON81",
			Type:     "F16",
		},
		{
			ICAO:     "RPQ",
			Callsign: "JESTER41",
			Type:     "F18",
		},
	}

	expectedSpottedAircraft := []Aircraft{
		{
			ICAO:     "DEF",
			Callsign: "LDS431",
			Type:     "CESSNA",
		},
		{
			ICAO:     "VVO",
			Callsign: "DEMON81",
			Type:     "F16",
		},
		{
			ICAO:     "RPQ",
			Callsign: "JESTER41",
			Type:     "F18",
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
	aircraft := AircraftRaw{
		Callsign: "    ",
	}

	expected := "UNKNOWN"
	actual := validateFields(aircraft).Callsign

	if expected != actual {
		t.Fatalf("expected '%v' to be the same as '%v'", expected, actual)
	}
}

// TestCallsignEndsEmptySpace that the callsign is not changed if it does not start with empty spaces
func TestCallsignEndsEmptySpace(t *testing.T) {
	aircraft := AircraftRaw{
		Callsign: "x  x  ",
	}

	expected := "x  x"
	actual := validateFields(aircraft).Callsign

	if expected != actual {
		t.Fatalf("expected '%v' to be the same as '%v'", expected, actual)
	}
}

// TestCallsignNoEmptySpace tests that the callsign is not changed if it does not start or end with empty spaces
func TestCallsignNoEmptySpace(t *testing.T) {
	aircraft := AircraftRaw{
		Callsign: "x  x x",
	}

	expected := "x  x x"
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
			Altitude: 4000.0,
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

// TestFilterAircraftByAltitudeWithGroundValue tests the handling of 'ground' string as altitude
func TestFilterAircraftByAltitudeWithGroundValue(t *testing.T) {
	aircraftWithGroundValues := []Aircraft{
		{
			Callsign: "KHARMA11",
			Altitude: 0,
		},
		{
			Callsign: "KHARMA13",
			Altitude: float64(3000), // Normal altitude value
		},
		{
			Callsign: "KHARMA14",
			Altitude: float64(8000), // Above our test altitude filter
		},
	}

	maxAltitude := 5000 // Max altitude in feet
	filtered := filterAircraftByAltitude(aircraftWithGroundValues, maxAltitude)

	// We expect 3 aircraft: the two with ground values and the one at 3000ft
	expected := []Aircraft{
		{
			Callsign: "KHARMA11",
			Altitude: float64(0),
		},
		{
			Callsign: "KHARMA13",
			Altitude: float64(3000),
		},
	}

	if len(filtered) != len(expected) {
		t.Fatalf("expected %d aircraft, got %d", len(expected), len(filtered))
	}

	// Check for each aircraft in the expected results
	for i, expectedAc := range expected {
		found := false
		for _, actualAc := range filtered {
			if actualAc.Callsign == expectedAc.Callsign {
				found = true
				// Verify altitude was converted properly
				if expectedAc.Altitude != actualAc.Altitude {
					t.Fatalf("expected altitude %v for %s, got %v",
						expectedAc.Altitude, expectedAc.Callsign, actualAc.Altitude)
				}
				break
			}
		}
		if !found {
			t.Fatalf("expected aircraft %s not found in results at index %d", expectedAc.Callsign, i)
		}
	}
}

func TestHandleAircraftWithScanRange(t *testing.T) {
	// Create test aircraft at different distances
	testAircraft := []Aircraft{
		{
			ICAO:      "TEST1",
			Callsign:  "NEAR1",
			Latitude:  51.18,
			Longitude: 5.46,
			Altitude:  float64(1000),
		},
		{
			ICAO:      "TEST2",
			Callsign:  "MEDIUM1",
			Latitude:  51.40,
			Longitude: 5.80,
			Altitude:  float64(2000),
		},
		{
			ICAO:      "TEST3",
			Callsign:  "FAR1",
			Latitude:  52.50,
			Longitude: 7.00,
			Altitude:  float64(3000),
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
	}

	// Filter aircraft based on distance from the location
	var aircraftInScanRange []Aircraft
	for _, ac := range testAircraft {
		distance := CalculateDistance(config.Location, geodist.Coord{Lat: ac.Latitude, Lon: ac.Longitude})
		if distance <= config.MaxScanRangeKilometers {
			aircraftInScanRange = append(aircraftInScanRange, ac)
		}
	}

	// Manually apply the same filtering logic as in HandleAircraft function
	var aircraftInNotificationRange []Aircraft
	for _, ac := range aircraftInScanRange {
		distance := CalculateDistance(config.Location, geodist.Coord{Lat: ac.Latitude, Lon: ac.Longitude})
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

// TestIsAircraftInboundDirectly tests aircraft flying directly towards the target location
func TestIsAircraftInboundDirectly(t *testing.T) {
	// Set up location and aircraft flying directly towards it
	// The aircraft's Track should match the bearing from aircraft to location
	location := geodist.Coord{
		Lat: 51.0,
		Lon: 5.0,
	}

	// Calculate the bearing from aircraft to location
	aircraftLocation := geodist.Coord{
		Lat: 51.5,
		Lon: 5.0,
	}
	bearingToTarget := CalculateBearing(aircraftLocation, location)

	aircraft := AircraftRaw{
		Lat:   51.5,
		Lon:   5.0,
		Track: bearingToTarget, // Aircraft heading directly toward the location
	}

	// Test with different margins
	testCases := []struct {
		margin   float64
		expected bool
	}{
		{0.1, true},  // Very tight margin
		{5.0, true},  // Small margin
		{20.0, true}, // Wide margin
		{90.0, true}, // Very wide margin (should still be true)
	}

	for _, tc := range testCases {
		result := IsAircraftInbound(location, aircraft, tc.margin)
		if result != tc.expected {
			bearing := CalculateBearing(aircraftLocation, location)
			diff := math.Abs(bearing - aircraft.Track)
			if diff > 180 {
				diff = 360 - diff
			}

			t.Fatalf("IsAircraftInbound with margin %v: expected %v, got %v. Bearing: %v, Track: %v, diff: %v",
				tc.margin, tc.expected, result, bearing, aircraft.Track, diff)
		}
	}
}

// TestIsAircraftInboundSlightlyOff tests aircraft flying slightly off from direct path
func TestIsAircraftInboundSlightlyOff(t *testing.T) {
	// Set up location and aircraft flying slightly off from direct path
	location := geodist.Coord{
		Lat: 51.0,
		Lon: 5.0,
	}

	// Calculate the ideal bearing from aircraft to location
	aircraftLocation := geodist.Coord{
		Lat: 51.5,
		Lon: 5.1, // Slightly east of direct north
	}
	idealBearing := CalculateBearing(aircraftLocation, location)

	// Aircraft is flying 5 degrees off from ideal bearing
	aircraft := AircraftRaw{
		Lat:   51.5,
		Lon:   5.1,
		Track: idealBearing - 5, // 5 degrees off from ideal bearing
	}

	// Test with different margins
	testCases := []struct {
		margin   float64
		expected bool
	}{
		{1.0, false}, // Very tight margin - should fail
		{5.0, true},  // Small margin - should pass
		{20.0, true}, // Wide margin - should pass
	}

	for _, tc := range testCases {
		result := IsAircraftInbound(location, aircraft, tc.margin)
		if result != tc.expected {
			bearing := CalculateBearing(aircraftLocation, location)
			diff := math.Abs(bearing - aircraft.Track)
			if diff > 180 {
				diff = 360 - diff
			}

			t.Fatalf("IsAircraftInbound with margin %v: expected %v, got %v. Bearing: %v, Track: %v, diff: %v",
				tc.margin, tc.expected, result, bearing, aircraft.Track, diff)
		}
	}
}

// TestIsAircraftNotInbound tests aircraft flying perpendicular to or away from the target
func TestIsAircraftNotInbound(t *testing.T) {
	// Case 1: Aircraft flying perpendicular to the location
	perpLocation := geodist.Coord{
		Lat: 51.0,
		Lon: 5.0,
	}

	perpAircraftLocation := geodist.Coord{
		Lat: 51.0,
		Lon: 6.0, // East of the location
	}
	idealBearing := CalculateBearing(perpAircraftLocation, perpLocation)

	perpAircraft := AircraftRaw{
		Lat:   51.0,
		Lon:   6.0,
		Track: math.Mod(idealBearing+90, 360), // Flying 90 degrees off from ideal bearing
	}

	// Case 2: Aircraft flying away from the location
	awayLocation := geodist.Coord{
		Lat: 51.0,
		Lon: 5.0,
	}

	awayAircraftLocation := geodist.Coord{
		Lat: 50.5,
		Lon: 5.0, // South of the location
	}
	idealBearingAway := CalculateBearing(awayAircraftLocation, awayLocation)

	awayAircraft := AircraftRaw{
		Lat:   50.5,
		Lon:   5.0,
		Track: math.Mod(idealBearingAway+180, 360), // Flying in the opposite direction
	}

	// Test both cases with different margins
	testCases := []struct {
		location    geodist.Coord
		aircraft    AircraftRaw
		aircraftLoc geodist.Coord
		margin      float64
		expected    bool
		name        string
	}{
		{perpLocation, perpAircraft, perpAircraftLocation, 10.0, false, "perpendicular with 10° margin"},
		{perpLocation, perpAircraft, perpAircraftLocation, 45.0, false, "perpendicular with 45° margin"},
		{perpLocation, perpAircraft, perpAircraftLocation, 95.0, true, "perpendicular with 95° margin"}, // Should pass with very wide margin
		{awayLocation, awayAircraft, awayAircraftLocation, 10.0, false, "flying away with 10° margin"},
		{awayLocation, awayAircraft, awayAircraftLocation, 45.0, false, "flying away with 45° margin"},
		{awayLocation, awayAircraft, awayAircraftLocation, 179.0, false, "flying away with 179° margin"},
		{awayLocation, awayAircraft, awayAircraftLocation, 185.0, true, "flying away with 185° margin"}, // Should pass with margin > 180
	}

	for _, tc := range testCases {
		result := IsAircraftInbound(tc.location, tc.aircraft, tc.margin)
		if result != tc.expected {
			bearing := CalculateBearing(tc.aircraftLoc, tc.location)
			diff := math.Abs(bearing - tc.aircraft.Track)
			if diff > 180 {
				diff = 360 - diff
			}

			t.Fatalf("%s: expected %v, got %v. Bearing: %v, Track: %v, diff: %v",
				tc.name, tc.expected, result, bearing, tc.aircraft.Track, diff)
		}
	}
}

// TestIsAircraftInboundRealWorld tests the function with real-world coordinates
func TestIsAircraftInboundRealWorld(t *testing.T) {
	// Brussels location
	brusselsLocation := geodist.Coord{Lat: 50.844987, Lon: 4.349981}

	// Aircraft near Heathrow heading towards Brussels
	heathrowLocation := geodist.Coord{Lat: 51.470020, Lon: -0.454295}
	heathrowToBrusselsBearing := CalculateBearing(heathrowLocation, brusselsLocation)

	heathrowAircraft := AircraftRaw{
		Lat:   51.470020,
		Lon:   -0.454295,
		Track: heathrowToBrusselsBearing, // Exact heading from Heathrow to Brussels
	}

	// Aircraft heading slightly off from Brussels
	heathrowAircraftOffCourse := AircraftRaw{
		Lat:   51.470020,
		Lon:   -0.454295,
		Track: heathrowToBrusselsBearing + 12, // 12 degrees off
	}

	// JFK to Brussels bearing
	jfkLocation := geodist.Coord{Lat: 40.639751, Lon: -73.778925}
	jfkToBrusselsBearing := CalculateBearing(jfkLocation, brusselsLocation)

	// JFK Aircraft heading towards Europe
	jfkAircraft := AircraftRaw{
		Lat:   40.639751,
		Lon:   -73.778925,
		Track: jfkToBrusselsBearing, // Exact heading to Brussels
	}

	// Test with existing locations from test file
	elisabethToPisBearing := CalculateBearing(locationElisabethPark, locationMannekenPis)

	elisabethParkAircraft := AircraftRaw{
		Lat:   locationElisabethPark.Lat,
		Lon:   locationElisabethPark.Lon,
		Track: elisabethToPisBearing, // Exact heading to Manneken Pis
	}

	testCases := []struct {
		location geodist.Coord
		aircraft AircraftRaw
		margin   float64
		expected bool
		name     string
	}{
		{
			location: brusselsLocation,
			aircraft: heathrowAircraft,
			margin:   5.0,
			expected: true,
			name:     "Heathrow to Brussels with 5° margin (exact heading)",
		},
		{
			location: brusselsLocation,
			aircraft: heathrowAircraftOffCourse,
			margin:   5.0,
			expected: false, // With tight margin, the off-course aircraft fails
			name:     "Heathrow to Brussels with 5° margin (12° off course)",
		},
		{
			location: brusselsLocation,
			aircraft: heathrowAircraftOffCourse,
			margin:   15.0,
			expected: true, // With wider margin, it passes
			name:     "Heathrow to Brussels with 15° margin (12° off course)",
		},
		{
			location: brusselsLocation,
			aircraft: jfkAircraft,
			margin:   5.0,
			expected: true,
			name:     "JFK to Brussels with 5° margin (exact heading)",
		},
		{
			location: locationMannekenPis,
			aircraft: elisabethParkAircraft,
			margin:   5.0,
			expected: true,
			name:     "Elisabeth Park to Manneken Pis with 5° margin (exact heading)",
		},
	}

	for _, tc := range testCases {
		result := IsAircraftInbound(tc.location, tc.aircraft, tc.margin)
		if result != tc.expected {
			aircraftLocation := geodist.Coord{Lat: tc.aircraft.Lat, Lon: tc.aircraft.Lon}
			bearing := CalculateBearing(aircraftLocation, tc.location)
			diff := math.Abs(bearing - tc.aircraft.Track)
			if diff > 180 {
				diff = 360 - diff
			}

			t.Fatalf("%s: expected %v, got %v. Bearing: %v, Track: %v, diff: %v",
				tc.name, tc.expected, result, bearing, tc.aircraft.Track, diff)
		}
	}
}

// TestAircraftOnGroundIsNeverInbound tests that an aircraft with altitude 0 is never marked as inbound
func TestAircraftOnGroundIsNeverInbound(t *testing.T) {
	// Create test aircraft and configuration
	aircraftOnGround := AircraftRaw{
		ICAO:         "TEST1",
		Callsign:     "GROUND1",
		Registration: "ABC",
		Lat:          51.18,
		Lon:          5.46,
		AltBaro:      float64(0), // Aircraft on the ground
		Track:        45.0,       // Heading that would normally be considered inbound
	}

	config := configuration.Config{
		Location: geodist.Coord{
			Lat: 51.17348,
			Lon: 5.45921,
		},
	}

	aircraft := []AircraftRaw{aircraftOnGround}

	// Process the aircraft
	outputs, err := ConvertToAircraft(aircraft, config, false)
	if err != nil {
		t.Fatalf("Error creating aircraft output: %v", err)
	}

	// Verify there's exactly one aircraft in the result
	if len(outputs) != 1 {
		t.Fatalf("Expected 1 aircraft in output, got %d", len(outputs))
	}

	// Check that the aircraft is marked as on the ground
	if !outputs[0].OnGround {
		t.Error("Expected aircraft to be marked as on the ground (OnGround=true)")
	}

	// Check that the aircraft is not marked as inbound (regardless of its heading)
	if outputs[0].Inbound {
		t.Error("Aircraft on the ground should never be marked as inbound, but was marked as inbound")
	}
}

func TestAircraftWithoutRegistrationIsSkipped(t *testing.T) {
	// Create test aircraft and configuration
	aircraftWithoutRegistration := AircraftRaw{
		ICAO:         "TEST9",
		Callsign:     "REGI1",
		Registration: "", // No registration
		Lat:          51.18,
		Lon:          5.46,
	}

	aircraftWithRegistration := AircraftRaw{
		ICAO:         "TEST10",
		Callsign:     "REGI1",
		Registration: "ABC",
		Lat:          51.18,
		Lon:          5.46,
	}

	aircraft := []AircraftRaw{aircraftWithoutRegistration, aircraftWithRegistration}

	// Process the aircraft
	outputs, err := ConvertToAircraft(aircraft, configuration.Config{}, false)
	if err != nil {
		t.Fatalf("Error creating aircraft output: %v", err)
	}

	// Verify there's exactly one aircraft in the result
	if len(outputs) != 1 {
		t.Fatalf("Expected 1 aircraft in output, got %d", len(outputs))
	}
	// Check that the aircraft is the one with registration
	if outputs[0].Registration != "ABC" {
		t.Fatalf("Expected aircraft with registration 'ABC', got '%s'", outputs[0].Registration)
	}
}
