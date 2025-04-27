package notification

import (
	"fmt"
	"jetspotter/internal/configuration"
	"jetspotter/internal/jetspotter"
)

type NtfyNotification struct {
	Message  string   `json:"message,omitempty"`
	Topic    string   `json:"topic,omitempty"`
	Title    string   `json:"title,omitempty"`
	Markdown bool     `json:"markdown,omitempty"`
	Tags     []string `json:"tags,omitempty"`
}

// SendNtfyMessage sends a ntfy message containing metadata of a list of aircraft
func SendNtfyMessage(aircraft []jetspotter.AircraftOutput, config configuration.Config) error {
	message, err := buildNtfyMessage(aircraft, config)
	if err != nil {
		return err
	}

	notification := Notification{
		Message: message,
		Type:    Ntfy,
		URL:     config.NtfyServer,
	}

	err = SendMessage(aircraft, notification)
	if err != nil {
		return err
	}

	return nil
}

func buildNtfyMessage(aircraft []jetspotter.AircraftOutput, config configuration.Config) (message NtfyNotification, err error) {
	message.Title = "An aircraft has been spotted!"
	message.Topic = config.NtfyTopic
	message.Tags = []string{"jetspotter"}
	message.Markdown = true

	for _, ac := range aircraft {
		message.Message += "\\=\\=\\=\\=\\=\\=\\=\\=\\=\\=\\=\\=\\=\\=\\=\\=\\=\\=\n"
		message.Message += fmt.Sprintf("Callsign: %s\n", formatCallsign(ac, Markdown))
		message.Message += fmt.Sprintf("Registration: %s\n", formatRegistration(ac, Markdown))
		message.Message += fmt.Sprintf("Country: %s\n", ac.Country)
		message.Message += fmt.Sprintf("Speed: %s\n", printSpeed(ac))
		message.Message += fmt.Sprintf("Altitude: %s\n", printAltitude(ac))
		message.Message += fmt.Sprintf("Distance: %s\n", printDistance(ac))
		message.Message += fmt.Sprintf("Bearing from location: %s\n", printBearingFromLocation(ac))
		message.Message += fmt.Sprintf("Bearing to location: %s\n", printBearingFromAircraft(ac))
		message.Message += fmt.Sprintf("Heading: %s\n", printHeading(ac))
		message.Message += fmt.Sprintf("Cloud coverage: %s\n", printCloudCoverage(ac))
		message.Message += fmt.Sprintf("Inbound: %s\n", getInboundStatus(ac))
		message.Message += fmt.Sprintf("Type: %s\n", ac.Type)
	}

	return message, nil
}
