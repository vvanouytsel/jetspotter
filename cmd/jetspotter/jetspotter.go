package main

import (
	"jetspotter/internal/configuration"
	"jetspotter/internal/jetspotter"
	"jetspotter/internal/notification"
	"log"
	"time"
)

func exitWithError(err error) {
	log.Fatalf("Something went wrong: %v\n", err)
}

func sendNotifications(aircraft []jetspotter.AircraftOutput, config configuration.Config) error {
	sortedAircraft := jetspotter.SortByDistance(aircraft)

	if len(aircraft) < 1 {
		log.Println("No new matching aircraft has been spotted.")
		return nil
	}

	// Terminal
	notification.SendTerminalMessage(sortedAircraft, config)

	// Slack
	if config.SlackWebHookURL != "" {
		err := notification.SendSlackMessage(sortedAircraft, config)
		if err != nil {
			return err
		}
	}

	// Discord
	if config.DiscordWebHookURL != "" {
		err := notification.SendDiscordMessage(sortedAircraft, config)
		if err != nil {
			return err
		}
	}

	return nil
}

func jetspotterHandler(alreadySpottedAircraft *[]jetspotter.Aircraft, config configuration.Config) {
	aircraft, err := jetspotter.HandleAircraft(alreadySpottedAircraft, config)
	if err != nil {
		exitWithError(err)
	}

	err = sendNotifications(aircraft, config)
	if err != nil {
		exitWithError(err)
	}
}

func main() {
	config, err := configuration.GetConfig()
	if err != nil {
		exitWithError(err)
	}

	var alreadySpottedAircraft []jetspotter.Aircraft
	for {

		if len(alreadySpottedAircraft) == 1 {
			log.Printf("%d aircraft is skipped, since it is already spotted.\n", len(alreadySpottedAircraft))
		}

		if len(alreadySpottedAircraft) > 0 {
			log.Printf("%d aircraft are skipped, since they are already spotted.\n", len(alreadySpottedAircraft))
		}
		jetspotterHandler(&alreadySpottedAircraft, config)
		time.Sleep(time.Duration(config.FetchInterval) * time.Second)
	}

}
