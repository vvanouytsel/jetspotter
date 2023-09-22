package jetspotter

import (
	"jetspotter/internal/aircraft"
	"reflect"
	"testing"
)

var (
	planes = []Aircraft{
		{
			Flight:    "APEX11",
			PlaneType: aircraft.F16.Identifier,
			Desc:      aircraft.F16.Description,
		},
		{
			Flight:    "XSG123",
			PlaneType: aircraft.B77L.Identifier,
			Desc:      aircraft.B77L.Description,
		},
	}
)

func TestFilterAircraftByTypeF16(t *testing.T) {
	expected := []Aircraft{
		{
			Flight:    "APEX11",
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
