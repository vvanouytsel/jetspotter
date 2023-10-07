package notification

import "testing"

func TestColorByAltitude(t *testing.T) {

	expected := darkBlue
	actual := getColorByAltitude(25000)

	if expected != actual {
		t.Fatalf("expected '%v' to be the same as '%v'", expected, actual)
	}
}
