package main

import (
	"jetspotter/internal/configuration"
	"jetspotter/internal/jetspotter"
	"jetspotter/internal/metrics"
	"jetspotter/internal/notification"
	"jetspotter/internal/version"
	"jetspotter/internal/web"
	"log"
	"strconv"
	"time"
)

// Global channel to indicate when the first data fetch is complete
var dataReadyChan = make(chan bool, 1)

func exitWithError(err error) {
	log.Fatalf("Something went wrong: %v\n", err)
}

func sendNotifications(aircraft []jetspotter.Aircraft, config configuration.Config) error {
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

func jetspotterHandler(alreadySpottedAircraft *[]jetspotter.Aircraft, config configuration.Config, isFirstRun bool) {
	aircraft, err := jetspotter.HandleAircraft(alreadySpottedAircraft, config)
	if err != nil {
		exitWithError(err)
	}

	err = sendNotifications(aircraft, config)
	if err != nil {
		exitWithError(err)
	}

	// If this is the first successful data fetch, signal that data is ready
	if isFirstRun {
		// Signal that data is ready (use non-blocking send)
		select {
		case dataReadyChan <- true:
			log.Println("Notified web UI that aircraft data is now available")
		default:
			// Channel already has a value or no receivers yet, that's fine
		}
	}
}

func HandleJetspotter(config configuration.Config) {
	if config.MaxScanRangeKilometers > config.MaxRangeKilometers {
		log.Printf("Scanning for aircraft within %d kilometers, sending notifications for those within %d kilometers: %s",
			config.MaxScanRangeKilometers, config.MaxRangeKilometers, config.AircraftTypes)
	} else {
		log.Printf("Spotting the following aircraft types within %d kilometers: %s",
			config.MaxRangeKilometers, config.AircraftTypes)
	}

	if config.MaxAltitudeFeet > 0 {
		log.Printf("Only showing aircraft at or below %d feet.", config.MaxAltitudeFeet)
	}

	var alreadySpottedAircraft []jetspotter.Aircraft
	isFirstRun := true

	for {
		jetspotterHandler(&alreadySpottedAircraft, config, isFirstRun)
		if isFirstRun {
			isFirstRun = false
		}
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

func HandleAPI(config configuration.Config) {
	jetspotter.SetupAPI(config.APIPort, config)
}

func HandleWebUI(config configuration.Config) {
	if !config.WebUIEnabled {
		return
	}

	go func() {
		// Convert port string to int
		port, err := strconv.Atoi(config.WebUIPort)
		if err != nil {
			log.Printf("Invalid web UI port: %s, using default: 8080", config.WebUIPort)
			port = 8080
		}

		// Set up API endpoint - we'll use localhost with the API port
		apiEndpoint := "http://localhost:" + config.APIPort

		// Configure and start the web server
		webConfig := web.Config{
			ListenPort:    port,
			APIEndpoint:   apiEndpoint,
			RefreshPeriod: time.Duration(config.FetchInterval) * time.Second,
			DataReadyChan: dataReadyChan, // Pass the channel to the web server
		}

		webServer := web.NewServer(webConfig)
		if err := webServer.Start(); err != nil {
			log.Printf("Web UI server error: %v", err)
		}
	}()
}

func main() {
	// Display a professional version banner at startup
	versionBanner := `
  ╔════════════════════════════════════════════════════════╗
  ║                      JETSPOTTER                        ║
  ╠════════════════════════════════════════════════════════╣
  ║  Version: %-10s                                   ║
  ║  Commit:  %-10s                                   ║
  ║  Built:   %-10s                                   ║
  ╚════════════════════════════════════════════════════════╝
`
	log.Printf(versionBanner, version.Version, version.Commit, version.BuildTime)

	config, err := configuration.GetConfig()
	if err != nil {
		exitWithError(err)
	}

	// Start services
	HandleMetrics(config)
	HandleAPI(config)
	HandleWebUI(config)

	// Start the main aircraft tracking loop
	HandleJetspotter(config)
}
