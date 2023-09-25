package main

import (
	"fmt"
	"jetspotter/internal/aircraft"
	"jetspotter/internal/jetspotter"
	notification "jetspotter/internal/notification"
	"os"
)

func exitWithError(err error) {
	fmt.Printf("Something went wrong: %v\n", err)
	os.Exit(1)
}

func main() {
	// TODO read this from environment variables if defined
	maxRangeKilometers := 50
	maxAmountAircraftSlackMessage := 20
	// aircraftType := aircraft.ALL.Identifier
	aircraftType := aircraft.F16.Identifier

	aircraft, err := jetspotter.GetFiltererdAircraftInRange(jetspotter.Bullseye, aircraftType, maxRangeKilometers)
	if err != nil {
		exitWithError(err)
	}

	jetspotter.PrintAircraft(aircraft)

	if len(aircraft) > 0 {
		err = notification.SendSlackMessage(aircraft, maxAmountAircraftSlackMessage)
		if err != nil {
			exitWithError(err)
		}
	}

}
