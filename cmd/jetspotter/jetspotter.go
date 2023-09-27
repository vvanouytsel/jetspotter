package main

import (
	"fmt"
	"jetspotter/internal/jetspotter"
	notification "jetspotter/internal/notification"
	"os"
	"strconv"
)

func exitWithError(err error) {
	fmt.Printf("Something went wrong: %v\n", err)
	os.Exit(1)
}

func main() {
	maxAmountAircraftSlackMessage := 8
	maxRangeKilometers, err := strconv.Atoi(jetspotter.GetEnvVariable("MAX_RANGE_KILOMETERS", "30"))
	if err != nil {
		exitWithError(err)
	}

	aircraftType := jetspotter.GetEnvVariable("AIRCRAFT_TYPE", "ALL")
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
