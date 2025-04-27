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
					Text: formatRegistration(ac, Slack),
				},
				{
					Type: "mrkdwn",
					Text: fmt.Sprintf("*Country:* %s", ac.Country),
				},
				{
					Type: "mrkdwn",
					Text: fmt.Sprintf("*Speed:* %s", printSpeed(ac)),
				},
				{
					Type: "mrkdwn",
					Text: fmt.Sprintf("*Altitude:* %s", printAltitude(ac)),
				},
				{
					Type: "mrkdwn",
					Text: fmt.Sprintf("*Distance:* %s", printDistance(ac)),
				},
				{
					Type: "mrkdwn",
					Text: fmt.Sprintf("*Bearing from location:* %s", printBearingFromLocation(ac)),
				},
				{
					Type: "mrkdwn",
					Text: fmt.Sprintf("*Heading:* %s", printHeading(ac)),
				},
				{
					Type: "mrkdwn",
					Text: fmt.Sprintf("*Bearing from aircraft:* %s", printBearingFromAircraft(ac)),
				},
				{
					Type: "mrkdwn",
					Text: fmt.Sprintf("*Cloud coverage:* %s", printCloudCoverage(ac)),
				},
				{
					Type: "mrkdwn",
					Text: fmt.Sprintf("*Inbound:* %s", getInboundStatus(ac)),
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
		Message: message,
		Type:    Slack,
		URL:     config.SlackWebHookURL,
	}

	err = SendMessage(aircraft, notification)
	if err != nil {
		return err
	}
	return nil
}
