package notification

import (
	"fmt"
	"jetspotter/internal/configuration"
	"jetspotter/internal/jetspotter"
	"log"
)

// FormatAircraft prints an Aircraft in a readable manner.
func FormatAircraft(aircraft jetspotter.AircraftOutput, config configuration.Config) string {

	return fmt.Sprintf("Callsign: %s\n"+
		"Description: %s\n"+
		"Type: %s\n"+
		"Tail number: %s\n"+
		"Country: %s\n"+
		"Altitude: %dft | %dm\n"+
		"Speed: %dkn | %dkm/h\n"+
		"Distance: %dkm\n"+
		"Cloud coverage: %d%%\n"+
		"Bearing from location: %.0f°\n"+
		"Bearing from aircraft: %.0f°\n"+
		"Heading: %.0f°\n"+
		"TrackerURL: %s\n"+
		"ImageURL: %s\n",

		aircraft.Callsign, aircraft.Description, aircraft.Type,
		aircraft.Registration, aircraft.Country, int(aircraft.Altitude), jetspotter.ConvertFeetToMeters(aircraft.Altitude),
		aircraft.Speed, jetspotter.ConvertKnotsToKilometersPerHour(aircraft.Speed),
		aircraft.Distance, aircraft.CloudCoverage, aircraft.BearingFromLocation,
		aircraft.BearingFromAircraft, aircraft.Heading, aircraft.TrackerURL, aircraft.ImageURL)
}

// SendTerminalMessage prints a list of Aircraft in a readable manner.
func SendTerminalMessage(aircraft []jetspotter.AircraftOutput, config configuration.Config) {
	log.Println("🛫 A jet has been spotted! 🛫")
	for _, ac := range aircraft {
		fmt.Println(FormatAircraft(ac, config))
	}
}
