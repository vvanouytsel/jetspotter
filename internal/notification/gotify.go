package notification

import (
	"fmt"
	"jetspotter/internal/configuration"
	"jetspotter/internal/jetspotter"

	"github.com/gotify/go-api-client/v2/models"
)

// SendGotifyMessage sends a discord message containing metadata of a list of aircraft
func SendGotifyMessage(aircraft []jetspotter.AircraftOutput, config configuration.Config) error {
	message, err := buildGotifyMessage(aircraft, config)
	if err != nil {
		return err
	}

	notification := Notification{
		Message: message,
		Type:    Gotify,
		URL:     fmt.Sprintf("%s/message?token=%s", config.GotifyURL, config.GotifyToken),
	}

	err = SendMessage(aircraft, notification)
	if err != nil {
		return err
	}

	return nil
}

func buildGotifyMessage(aircraft []jetspotter.AircraftOutput, config configuration.Config) (message models.MessageExternal, err error) {
	message.Title = "An aircraft has been spotted!"
	message.Extras = map[string]interface{}{
		"client::display": map[string]interface{}{
			"contentType": "text/markdown",
		},
	}

	for _, ac := range aircraft {
		message.Message += "==================\n\n"
		message.Message += fmt.Sprintf("**Callsign**: [%s](%s)\n\n", ac.Callsign, ac.TrackerURL)
		message.Message += fmt.Sprintf("**Registration**: %s\n\n", formatRegistration(ac, Markdown))
		message.Message += fmt.Sprintf("**Speed:** %s\n\n", printSpeed(ac))
		message.Message += fmt.Sprintf("**Altitude**: %s\n\n", printAltitude(ac))
		message.Message += fmt.Sprintf("**Distance:** %s\n\n", printDistance(ac))
		message.Message += fmt.Sprintf("**Bearing from location:** %s\n\n", printBearingFromLocation(ac))
		message.Message += fmt.Sprintf("**Bearing to location:** %s\n\n", printBearingFromAircraft(ac))
		message.Message += fmt.Sprintf("**Heading:** %s\n\n", printHeading(ac))
		message.Message += fmt.Sprintf("**Cloud coverage:** %s\n\n", printCloudCoverage(ac))
		message.Message += fmt.Sprintf("**Type:** %s\n\n", ac.Type)
	}

	return message, nil
}
