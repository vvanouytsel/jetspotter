package notification

import (
	"fmt"
	"jetspotter/internal/configuration"
	"jetspotter/internal/jetspotter"
)

type NtfyAction struct {
	Action string `json:"action"`
	Clear  bool   `json:"clear,omitempty"`
	Label  string `json:"label"`
	URL    string `json:"url"`
}

type NtfyNotification struct {
	Actions  []NtfyAction `json:"actions,omitempty"`
	Message  string       `json:"message,omitempty"`
	Topic    string       `json:"topic,omitempty"`
	Title    string       `json:"title,omitempty"`
	Markdown bool         `json:"markdown,omitempty"`
	Tags     []string     `json:"tags,omitempty"`
}

// SendNtfyMessage sends a ntfy message containing metadata of aircraft
// Each aircraft will have its own separate notification
func SendNtfyMessage(aircraft []jetspotter.Aircraft, config configuration.Config) error {
	// Send a separate message for each aircraft
	for _, ac := range aircraft {
		// Build a message for a single aircraft
		singleAircraftMessage, err := buildNtfyMessage(ac, config)
		if err != nil {
			return err
		}

		notification := Notification{
			Message: singleAircraftMessage,
			Type:    Ntfy,
			URL:     config.NtfyServer,
		}

		err = SendMessage(notification)
		if err != nil {
			return err
		}
	}

	return nil
}

// Constructor function for NtfyAction thats sets default value for Action and Clear
func AddNtfyAction(label, url string) NtfyAction {
	return NtfyAction{
		Action: "view",
		Clear:  true,
		Label:  label,
		URL:    url,
	}
}

func buildNtfyMessage(aircraft jetspotter.Aircraft, config configuration.Config) (message NtfyNotification, err error) {
	message.Title = "An aircraft has been spotted!"
	message.Topic = config.NtfyTopic
	message.Tags = []string{"jetspotter"}
	message.Markdown = true

	message.Message += fmt.Sprintf("Callsign:               %s\n", formatCallsign(aircraft, Markdown))
	message.Message += fmt.Sprintf("Registration:           %s\n", formatRegistration(aircraft, Markdown))
	message.Message += fmt.Sprintf("Country:                %s\n", aircraft.Country)
	message.Message += fmt.Sprintf("Speed:                  %s\n", printSpeed(aircraft))
	message.Message += fmt.Sprintf("Altitude:               %s\n", printAltitude(aircraft))
	message.Message += fmt.Sprintf("Distance:               %s\n", printDistance(aircraft))
	message.Message += fmt.Sprintf("Bearing from location:  %s\n", printBearingFromLocation(aircraft))
	message.Message += fmt.Sprintf("Bearing to location:    %s\n", printBearingFromAircraft(aircraft))
	message.Message += fmt.Sprintf("Heading:                %s\n", printHeading(aircraft))
	message.Message += fmt.Sprintf("Cloud coverage:         %s\n", printCloudCoverage(aircraft))
	message.Message += fmt.Sprintf("Inbound:                %s\n", getInboundStatus(aircraft))
	message.Message += fmt.Sprintf("Type:                   %s\n", aircraft.Type)
	message.Message += fmt.Sprintf("Origin:                 %s\n", printOriginName(aircraft))
	message.Message += fmt.Sprintf("Destination:            %s\n", printDestinationName(aircraft))
	message.Message += fmt.Sprintf("Airline:                %s\n", printAirlineName(aircraft))
	message.Message += fmt.Sprintf("ImageURL:               %s\n", aircraft.ImageURL)
	// Add Ntfy Actions
	message.Actions = []NtfyAction{
		AddNtfyAction("Track Aircraft", aircraft.TrackerURL),
	}

	return message, nil
}
