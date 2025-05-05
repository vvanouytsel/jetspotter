package notification

import (
	"bytes"
	"encoding/json"
	"fmt"
	"jetspotter/internal/jetspotter"
	"log"
	"net/http"
)

// Notification is a representation of the notfication that has to be sent
type Notification struct {
	Message interface{}
	Type    string
	URL     string
}

const (
	// Discord indicates the discord platform
	Discord = "Discord"
	// Slack indicates the slack platform
	Slack = "Slack"
	// Gotify indicates the gotify platform
	Gotify = "Gotify"
	// Ntfy indicates the ntfy platform
	Ntfy = "Ntfy"
	// Markdown indicates markdown markup language
	Markdown = "Markdown"
)

// Format the callsign whether to display a hyperlink or not
func formatCallsign(ac jetspotter.Aircraft, notificationType string) string {
	if notificationType == Markdown {
		return fmt.Sprintf("[%s](%s)", ac.Callsign, ac.TrackerURL)
	}

	if notificationType == Slack {
		return fmt.Sprintf("*Callsign:* <%s|%s>", ac.TrackerURL, ac.Callsign)
	}

	return ac.Callsign
}

// Format whether to display a hyperlink for the registration or not
func formatRegistration(ac jetspotter.Aircraft, notificationType string) string {
	if notificationType == Markdown {
		if ac.ImageURL == "" {
			return ac.Registration
		}

		return fmt.Sprintf("[%s](%s)", ac.Registration, ac.ImageURL)
	}

	if notificationType == Slack {
		if ac.ImageURL == "" {
			return fmt.Sprintf("*Registration:* %s", ac.Registration)
		}

		return fmt.Sprintf("*Registration:* <%s|%s>", ac.ImageURL, ac.Registration)
	}

	if ac.ImageURL == "" {
		return ac.Registration
	}

	return ac.Registration
}

// SendMessage sends a message to a notification platform
func SendMessage(notification Notification) error {
	data, err := json.Marshal(notification.Message)
	if err != nil {
		return err
	}

	resp, err := http.Post(notification.URL, "application/json",
		bytes.NewReader(data))
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		log.Printf("%s\n", string(data))
		return fmt.Errorf("received status code %v", resp.StatusCode)
	}

	log.Printf("A %s notification has been sent!\n", notification.Type)
	return nil
}

func printSpeed(ac jetspotter.Aircraft) string {
	return fmt.Sprintf("%dkn | %dkm/h", ac.Speed, jetspotter.ConvertKnotsToKilometersPerHour(ac.Speed))
}

func printAltitude(ac jetspotter.Aircraft) string {
	if ac.OnGround {
		return "On ground"
	}
	return fmt.Sprintf("%vft | %dm", ac.Altitude, jetspotter.ConvertFeetToMeters(ac.Altitude))
}

func printDistance(ac jetspotter.Aircraft) string {
	return fmt.Sprintf("%dkm", ac.Distance)
}

func printBearingFromLocation(ac jetspotter.Aircraft) string {
	return fmt.Sprintf("%.0f°", ac.BearingFromLocation)
}

func printHeading(ac jetspotter.Aircraft) string {
	return fmt.Sprintf("%.0f°", ac.Heading)
}

func printBearingFromAircraft(ac jetspotter.Aircraft) string {
	return fmt.Sprintf("%.0f°", ac.BearingFromAircraft)
}

func printCloudCoverage(ac jetspotter.Aircraft) string {
	return fmt.Sprintf("%d%%", ac.CloudCoverage)
}

func getInboundStatus(ac jetspotter.Aircraft) string {
	if ac.Inbound {
		return "Yes"
	}
	return "No"
}

func printOriginName(ac jetspotter.Aircraft) string {
	if ac.Origin.Name == "" {
		return "N/A"
	}
	return ac.Origin.Name
}

func printDestinationName(ac jetspotter.Aircraft) string {
	if ac.Destination.Name == "" {
		return "N/A"
	}
	return ac.Destination.Name
}

func printAirlineName(ac jetspotter.Aircraft) string {
	if ac.Airline.Name == "" {
		return "N/A"
	}
	return ac.Airline.Name
}
