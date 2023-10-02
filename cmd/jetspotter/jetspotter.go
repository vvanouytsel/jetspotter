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

func sendNotifications(aircraft []jetspotter.AircraftOutput, config configuration.Config) error {
	// CLI
	jetspotter.PrintAircraft(aircraft, config)

	// Slack
	if config.SlackWebHookURL != "" {
		err := notification.SendSlackMessage(aircraft, config)
		if err != nil {
			return err
		}
	}

	return nil
}

func main() {
	config, err := configuration.GetConfig()
	if err != nil {
		exitWithError(err)
	}

	aircraft, err := jetspotter.GetFiltererdAircraftInRange(config)
	if err != nil {
		exitWithError(err)
	}

	err = sendNotifications(aircraft, config)
	if err != nil {
		exitWithError(err)
	}

}
