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
	Message    interface{}
	Type       string
	WebHookURL string
}

const (
	// Discord indicates the discord platform
	Discord = "Discord"
	// Slack indicates the slack platform
	Slack = "Slack"
)

// SendMessage sends a message to a notification platform
func SendMessage(aircraft []jetspotter.AircraftOutput, notification Notification) error {

	data, err := json.Marshal(notification.Message)
	if err != nil {
		return err
	}

	resp, err := http.Post(notification.WebHookURL, "application/json",
		bytes.NewBuffer(data))
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		log.Printf("%s\n", string(data))
		return fmt.Errorf(fmt.Sprintf("Received status code %v", resp.StatusCode))
	}

	log.Printf("A %s notification has been sent!\n", notification.Type)
	return nil
}
