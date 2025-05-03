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

// Constructor function for NtfyAction thats sets default value for Action and Clear
func AddNtfyAction(label, url string) NtfyAction {
	return NtfyAction{
		Action: "view",
		Clear:  true,
		Label:  label,
		URL:    url,
	}
}

func buildNtfyMessage(aircraft []jetspotter.AircraftOutput, config configuration.Config) (message NtfyNotification, err error) {
	message.Title = "An aircraft has been spotted!"
	message.Topic = config.NtfyTopic
	message.Tags = []string{"jetspotter"}
	message.Markdown = true

	for _, ac := range aircraft {
		message.Message += fmt.Sprintf("Callsign:               %s\n", formatCallsign(ac, Markdown))
		message.Message += fmt.Sprintf("Registration:           %s\n", formatRegistration(ac, Markdown))
		message.Message += fmt.Sprintf("Country:                %s\n", ac.Country)
		message.Message += fmt.Sprintf("Speed:                  %s\n", printSpeed(ac))
		message.Message += fmt.Sprintf("Altitude:               %s\n", printAltitude(ac))
		message.Message += fmt.Sprintf("Distance:               %s\n", printDistance(ac))
		message.Message += fmt.Sprintf("Bearing from location:  %s\n", printBearingFromLocation(ac))
		message.Message += fmt.Sprintf("Bearing to location:    %s\n", printBearingFromAircraft(ac))
		message.Message += fmt.Sprintf("Heading:                %s\n", printHeading(ac))
		message.Message += fmt.Sprintf("Cloud coverage:         %s\n", printCloudCoverage(ac))
		message.Message += fmt.Sprintf("Inbound:                %s\n", getInboundStatus(ac))
		message.Message += fmt.Sprintf("Type:                   %s\n", ac.Type)
		message.Message += fmt.Sprintf("Origin:                 %s\n", printOriginName(ac))
		message.Message += fmt.Sprintf("Destination:            %s\n", printDestinationName(ac))
		message.Message += fmt.Sprintf("Airline:                %s\n", printAirlineName(ac))
		message.Message += fmt.Sprintf("ImageURL:               %s\n", ac.ImageURL)
		// Add Ntfy Actions
		message.Actions = []NtfyAction{
			AddNtfyAction("Track Aircraft", ac.TrackerURL),
		}
	}

	return message, nil
}
