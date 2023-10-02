package jetspotter

import (
	"jetspotter/internal/aircraft"
	"jetspotter/internal/configuration"
	"reflect"
	"testing"

	"github.com/jftuga/geodist"
)

var (
	planes = []AircraftOutput{
		{
			Callsign:    "APEX11",
			Type:        aircraft.F16.Identifier,
			Description: aircraft.F16.Description,
		},
		{
			Callsign:    "APEX12",
			Type:        aircraft.F16.Identifier,
			Description: aircraft.F16.Description,
		},
		{
			Callsign:    "XSG123",
			Type:        aircraft.B77L.Identifier,
			Description: aircraft.B77L.Description,
		},
		{
			Callsign:    "GRZLY11",
			Type:        aircraft.A400.Identifier,
			Description: aircraft.A400.Description,
		},
	}
)

func TestFilterAircraftByTypeF16(t *testing.T) {
	expected := []AircraftOutput{
		{
			Callsign:    "APEX11",
			Type:        aircraft.F16.Identifier,
			Description: aircraft.F16.Description,
		},
		{
			Callsign:    "APEX12",
			Type:        aircraft.F16.Identifier,
			Description: aircraft.F16.Description,
		},
	}

	config := configuration.Config{
		AircraftTypes: []string{aircraft.F16.Identifier},
	}
	actual := filterAircraftByTypes(planes, config)

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("expected '%v' to be the same as '%v'", expected, actual)
	}
}

func TestFilterAircraftByTypeALL(t *testing.T) {
	config := configuration.Config{
		AircraftTypes: []string{aircraft.ALL.Identifier},
	}
	expected := planes
	actual := filterAircraftByTypes(planes, config)

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("expected '%v' to be the same as '%v'", expected, actual)
	}
}

func TestCalculateDistance(t *testing.T) {
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

func TestFilterAircraftByTypes(t *testing.T) {
	expected := []AircraftOutput{
		{
			Callsign:    "APEX11",
			Type:        aircraft.F16.Identifier,
			Description: aircraft.F16.Description,
		},
		{
			Callsign:    "APEX12",
			Type:        aircraft.F16.Identifier,
			Description: aircraft.F16.Description,
		},
		{
			Callsign:    "GRZLY11",
			Type:        aircraft.A400.Identifier,
			Description: aircraft.A400.Description,
		},
	}

	config := configuration.Config{
		AircraftTypes: []string{aircraft.F16.Identifier, aircraft.A400.Identifier},
	}
	actual := filterAircraftByTypes(planes, config)

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
