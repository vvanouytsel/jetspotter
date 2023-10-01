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

	types := []string{aircraft.F16.Identifier}
	actual := filterAircraftByTypes(planes, types)

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("expected '%v' to be the same as '%v'", expected, actual)
	}
}

func TestFilterAircraftByTypeALL(t *testing.T) {
	types := []string{"ALL"}
	expected := planes
	actual := filterAircraftByTypes(planes, types)

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

	types := []string{aircraft.F16.Identifier, aircraft.A400.Identifier}
	actual := filterAircraftByTypes(planes, types)

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("expected '%v' to be the same as '%v'", expected, actual)
	}
}
