package main

import (
	"jetspotter/internal/configuration"
	"jetspotter/internal/jetspotter"
	"jetspotter/internal/metrics"
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
		log.Println("No new matching aircraft have been spotted.")
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

	// Gotify
	if config.GotifyURL != "" && config.GotifyToken != "" {
		err := notification.SendGotifyMessage(sortedAircraft, config)
		if err != nil {
			return err
		}
	}

	// Ntfy
	if config.NtfyTopic != "" {
		err := notification.SendNtfyMessage(sortedAircraft, config)
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

func HandleJetspotter(config configuration.Config) {
	log.Printf("Spotting the following aircraft types within %d kilometers: %s", config.MaxRangeKilometers, config.AircraftTypes)

	var alreadySpottedAircraft []jetspotter.Aircraft
	for {
		jetspotterHandler(&alreadySpottedAircraft, config)
		time.Sleep(time.Duration(config.FetchInterval) * time.Second)
	}
}

func HandleMetrics(config configuration.Config) {
	go func() {
		err := metrics.HandleMetrics(config)
		if err != nil {
			exitWithError(err)
		}
	}()
}

func main() {
	config, err := configuration.GetConfig()
	if err != nil {
		exitWithError(err)
	}
	HandleMetrics(config)
	HandleJetspotter(config)
}
