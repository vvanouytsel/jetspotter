package notification

import (
	"fmt"
	"jetspotter/internal/configuration"
	"jetspotter/internal/jetspotter"
	"log"
)

// FormatAircraft prints an Aircraft in a readable manner.
func FormatAircraft(aircraft jetspotter.Aircraft, config configuration.Config) string {
	return fmt.Sprintf("Callsign: %s\n"+
		"Description: %s\n"+
		"Type: %s\n"+
		"Tail number: %s\n"+
		"Country: %s\n"+
		"Altitude: %s\n"+
		"Speed: %s\n"+
		"Distance: %s\n"+
		"Cloud coverage: %s\n"+
		"Bearing from location: %s\n"+
		"Bearing from aircraft: %s\n"+
		"Heading: %s\n"+
		"Inbound: %s\n"+
		"Origin: %s\n"+
		"Destination: %s\n"+
		"Airline: %s\n"+
		"TrackerURL: %s\n"+
		"ImageURL: %s\n",

		aircraft.Callsign, aircraft.Description, aircraft.Type,
		aircraft.Registration, aircraft.Country, printAltitude(aircraft),
		printSpeed(aircraft), printDistance(aircraft), printCloudCoverage(aircraft),
		printBearingFromLocation(aircraft), printBearingFromAircraft(aircraft),
		printHeading(aircraft), getInboundStatus(aircraft), printOriginName(aircraft),
		printDestinationName(aircraft), printAirlineName(aircraft), aircraft.TrackerURL, aircraft.ImageURL)
}

// SendTerminalMessage prints a list of Aircraft in a readable manner.
func SendTerminalMessage(aircraft []jetspotter.Aircraft, config configuration.Config) {
	log.Println("🛫 A jet has been spotted! 🛫")
	for _, ac := range aircraft {
		fmt.Println(FormatAircraft(ac, config))
	}
}
