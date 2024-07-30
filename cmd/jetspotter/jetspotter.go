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

func sendNotifications(aircraft []jetspotter.AircraftOutput) error {
	sortedAircraft := jetspotter.SortByDistance(aircraft)

	if len(aircraft) < 1 {
		log.Println("No new matching aircraft have been spotted.")
		return nil
	}

	// Terminal
	notification.SendTerminalMessage(sortedAircraft)

	// Slack
	if configuration.SlackWebhookURL != "" {
		err := notification.SendSlackMessage(sortedAircraft)
		if err != nil {
			return err
		}
	}

	// Discord
	if configuration.DiscordWebhookURL != "" {
		err := notification.SendDiscordMessage(sortedAircraft)
		if err != nil {
			return err
		}
	}

	// Gotify
	if configuration.GotifyURL != "" && configuration.GotifyToken != "" {
		err := notification.SendGotifyMessage(sortedAircraft)
		if err != nil {
			return err
		}
	}

	// Ntfy
	if configuration.NtfyTopic != "" {
		err := notification.SendNtfyMessage(sortedAircraft)
		if err != nil {
			return err
		}
	}

	return nil
}

func jetspotterHandler(alreadySpottedAircraft *[]jetspotter.Aircraft) {
	aircraft, err := jetspotter.HandleAircraft(alreadySpottedAircraft)
	if err != nil {
		exitWithError(err)
	}

	err = sendNotifications(aircraft)
	if err != nil {
		exitWithError(err)
	}
}

func HandleJetspotter() {
	log.Printf("Spotting the following aircraft types within %d kilometers: %s", configuration.MaxRangeKilometers, configuration.AircraftTypes)

	var alreadySpottedAircraft []jetspotter.Aircraft
	for {
		jetspotterHandler(&alreadySpottedAircraft)
		time.Sleep(time.Duration(configuration.FetchInterval) * time.Second)
	}
}

func HandleMetrics() {
	go func() {
		err := metrics.HandleMetrics()
		if err != nil {
			exitWithError(err)
		}
	}()
}

func main() {
	configuration.GetConfig()

	HandleMetrics()
	HandleJetspotter()
}
