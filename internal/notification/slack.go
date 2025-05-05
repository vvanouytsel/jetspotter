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

func buildSlackMessage(aircraft []jetspotter.Aircraft) (SlackMessage, error) {

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

	for _, ac := range aircraft {
		// First section block with first 8 fields
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
			},
		})

		// Second section block with remaining 7 fields
		blocks = append(blocks, Block{
			Type: "section",
			Fields: []Field{
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
				{
					Type: "mrkdwn",
					Text: fmt.Sprintf("*Origin:* %s", printOriginName(ac)),
				},
				{
					Type: "mrkdwn",
					Text: fmt.Sprintf("*Destination:* %s", printDestinationName(ac)),
				},
				{
					Type: "mrkdwn",
					Text: fmt.Sprintf("*Airline:* %s", printAirlineName(ac)),
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
func SendSlackMessage(aircraft []jetspotter.Aircraft, config configuration.Config) error {
	// Split aircraft into chunks to stay within Slack's block limit (max 50 blocks per message)
	// Each aircraft uses approximately 4 blocks (2 sections + image + divider)
	const maxAircraftPerMessage = 10

	for i := 0; i < len(aircraft); i += maxAircraftPerMessage {
		end := i + maxAircraftPerMessage
		if end > len(aircraft) {
			end = len(aircraft)
		}

		chunk := aircraft[i:end]
		message, err := buildSlackMessage(chunk)
		if err != nil {
			return err
		}

		notification := Notification{
			Message: message,
			Type:    Slack,
			URL:     config.SlackWebHookURL,
		}

		err = SendMessage(notification)
		if err != nil {
			return fmt.Errorf("failed to send slack message: %w", err)
		}
	}

	return nil
}
