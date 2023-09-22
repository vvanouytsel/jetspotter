package main

import (
	"fmt"
	"jetspotter/internal/aircraft"
	"jetspotter/internal/jetspotter"
	"os"
)

func exitWithError(err error) {
	fmt.Printf("Something went wrong: %v\n", err)
	os.Exit(1)
}

func main() {
	// TODO read this from environment variables if defined
	maxRangeKilometers := 50
	aircraftType := aircraft.ALL.Identifier
	// aircraftType := aircraft.F16.Identifier

	aircraft, err := jetspotter.GetAircraftTypeInRange(jetspotter.BullsEye, aircraftType, maxRangeKilometers)
	if err != nil {
		exitWithError(err)
	}

	jetspotter.PrintAircraft(aircraft)

	// if len(vipers) > 0 {
	// 	fmt.Println("****** SPOTTED AN F16, SHOULD SEND SLACK MESSAGE! ******")
	// 	for _, viper := range vipers {
	// 		// Calculate distance based on lat lon
	// 		fmt.Printf("CALLSIGN: %s\nSQUAWK: %s\nTAIL: %s\nLAT: %v\nLON: %v\n\n",
	// 			viper.Flight, viper.Squawk, viper.TailNumber, viper.Lat, viper.Lon)

	// 	}
	// }
}
