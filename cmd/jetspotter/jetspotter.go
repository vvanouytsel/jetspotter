package main

import (
	"fmt"
	"jetspotter/internal/configuration"
	"jetspotter/internal/jetspotter"
	notification "jetspotter/internal/notification"
	"os"
)

func exitWithError(err error) {
	fmt.Printf("Something went wrong: %v\n", err)
	os.Exit(1)
}

func main() {
	config, err := configuration.GetConfig()
	if err != nil {
		exitWithError(err)
	}

	aircraft, err := jetspotter.GetFiltererdAircraftInRange(config.Location, config.AircraftType, config.MaxRangeKilometers)
	if err != nil {
		exitWithError(err)
	}

	jetspotter.PrintAircraft(aircraft, config)

	if len(aircraft) > 0 {
		err = notification.SendSlackMessage(aircraft, config)
		if err != nil {
			exitWithError(err)
		}
	}

}
