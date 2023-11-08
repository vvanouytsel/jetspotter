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
	planes = []Aircraft{
		{
			Callsign:  "APEX11",
			PlaneType: aircraft.F16.Identifier,
			Desc:      aircraft.F16.Description,
		},
		{
			Callsign:  "APEX12",
			PlaneType: aircraft.F16.Identifier,
			Desc:      aircraft.F16.Description,
		},
		{
			Callsign:  "XSG123",
			PlaneType: aircraft.B77L.Identifier,
			Desc:      aircraft.B77L.Description,
		},
		{
			Callsign:  "GRZLY11",
			PlaneType: aircraft.A400.Identifier,
			Desc:      aircraft.A400.Description,
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
		},
		{
			Callsign:  "APEX12",
			PlaneType: aircraft.F16.Identifier,
			Desc:      aircraft.F16.Description,
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
	expected := planes
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
		},
		{
			Callsign:  "APEX12",
			PlaneType: aircraft.F16.Identifier,
			Desc:      aircraft.F16.Description,
		},
		{
			Callsign:  "GRZLY11",
			PlaneType: aircraft.A400.Identifier,
			Desc:      aircraft.A400.Description,
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
