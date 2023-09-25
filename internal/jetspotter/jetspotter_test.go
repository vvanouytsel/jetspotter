package jetspotter

import (
	"jetspotter/internal/aircraft"
	"reflect"
	"testing"

	"github.com/jftuga/geodist"
)

var (
	planes = []Aircraft{
		{
			Callsign:  "APEX11",
			PlaneType: aircraft.F16.Identifier,
			Desc:      aircraft.F16.Description,
		},
		{
			Callsign:  "XSG123",
			PlaneType: aircraft.B77L.Identifier,
			Desc:      aircraft.B77L.Description,
		},
	}
)

func TestFilterAircraftByTypeF16(t *testing.T) {
	expected := []Aircraft{
		{
			Callsign:  "APEX11",
			PlaneType: aircraft.F16.Identifier,
			Desc:      aircraft.F16.Description,
		}}

	actual := filterAircraftByType(planes, aircraft.F16.Identifier)

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("expected '%v' to be the same as '%v'", expected, actual)
	}
}

func TestFilterAircraftByTypeALL(t *testing.T) {
	expected := planes
	actual := filterAircraftByType(planes, aircraft.ALL.Identifier)

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
