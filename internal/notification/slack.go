package notification

import (
	"fmt"
	"jetspotter/internal/configuration"
	"jetspotter/internal/jetspotter"
)

// SlackMessage is used to create a slack message
type SlackMessage struct {
	Blocks []Block `json:"blocks"`
}

// Block is a single block in a slack message
type Block struct {
	Type     string  `json:"type,omitempty"`
	Fields   []Field `json:"fields,omitempty"`
	Title    *Title  `json:"title,omitempty"`
	ImageURL string  `json:"image_url,omitempty"`
	AltText  string  `json:"alt_text,omitempty"`
}

// Title is a tile in a block of a slack message
type Title struct {
	Type  string `json:"type,omitempty"`
	Text  string `json:"text,omitempty"`
	Emoji bool   `json:"emoji,omitempty"`
}

// Field is a field in a block of a slack message
type Field struct {
	Type string `json:"type,omitempty"`
	Text string `json:"text,omitempty"`
}

func buildSlackMessage(aircraft []jetspotter.AircraftOutput, config configuration.Config) (SlackMessage, error) {

	var blocks []Block
	blocks = append(blocks, Block{
		Type: "section",
		Fields: []Field{
			{
				Type: "mrkdwn",
				Text: ":airplane: A jet has been spotted! :airplane:",
			},
		},
	})

	for i, ac := range aircraft {
		if i > config.MaxAircraftSlackMessage {
			break
		}

		blocks = append(blocks, Block{
			Type: "section",
			Fields: []Field{
				{
					Type: "mrkdwn",
					Text: fmt.Sprintf("*Callsign:* <%s|%s>", ac.TrackerURL, ac.Callsign),
				},
				{
					Type: "mrkdwn",
					Text: fmt.Sprintf("*Registration:* <%s|%s>", ac.ImageURL, ac.Registration),
				},
				{
					Type: "mrkdwn",
					Text: fmt.Sprintf("*Speed:* %dkn | %dkm/h", ac.Speed, jetspotter.ConvertKnotsToKilometersPerHour(ac.Speed)),
				},
				{
					Type: "mrkdwn",
					Text: fmt.Sprintf("*Altitude:* %vft | %dm", ac.Altitude, jetspotter.ConvertFeetToMeters(ac.Altitude)),
				},
				{
					Type: "mrkdwn",
					Text: fmt.Sprintf("*Distance:* %dkm", ac.Distance),
				},
				{
					Type: "mrkdwn",
					Text: fmt.Sprintf("*Bearing from location:* %.0f°", ac.BearingFromLocation),
				},
				{
					Type: "mrkdwn",
					Text: fmt.Sprintf("*Heading:* %.0f°", ac.Heading),
				},
				{
					Type: "mrkdwn",
					Text: fmt.Sprintf("*Bearing from aircraft:* %.0f°", ac.BearingFromAircraft),
				},
				{
					Type: "mrkdwn",
					Text: fmt.Sprintf("*Cloud coverage:* %d%%", ac.CloudCoverage),
				},
				{
					Type: "mrkdwn",
					Text: fmt.Sprintf("*Type:* %s", ac.Description),
				},
			},
		})

		imageURL := ac.ImageThumbnailURL
		if imageURL != "" {
			blocks = append(blocks,
				Block{
					Type: "image",
					Title: &Title{
						Type:  "plain_text",
						Text:  fmt.Sprintf("%s - %s", ac.Description, ac.Registration),
						Emoji: true,
					},
					ImageURL: imageURL,
					AltText:  fmt.Sprintf("%s with registration number %s", ac.Description, ac.Registration),
				})
		}

		blocks = append(blocks,
			Block{
				Type: "divider",
			},
		)
	}

	slackMessage := SlackMessage{Blocks: blocks}
	return slackMessage, nil
}

// SendSlackMessage sends a slack message containing metadata of a list of aircraft
func SendSlackMessage(aircraft []jetspotter.AircraftOutput, config configuration.Config) error {
	message, err := buildSlackMessage(aircraft, config)
	if err != nil {
		return err
	}

	notification := Notification{
		Message:    message,
		Type:       Slack,
		WebHookURL: config.SlackWebHookURL,
	}

	err = SendMessage(aircraft, notification)
	if err != nil {
		return err
	}
	return nil
}
