package notification

import (
	"bytes"
	"encoding/json"
	"fmt"
	"jetspotter/internal/jetspotter"
	"net/http"
	"os"
)

type SlackMessage struct {
	Blocks []Block `json:"blocks"`
}

type Block struct {
	Type   string  `json:"type"`
	Fields []Field `json:"fields,omitempty"`
}

type Field struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

func buildMessage([]jetspotter.Aircraft) SlackMessage {
	// TODO actually build it properly instead of hardcoding
	slackMessage := SlackMessage{
		Blocks: []Block{
			{
				Type: "section",
				Fields: []Field{
					{
						Type: "mrkdwn",
						Text: "A jet has been spotted:\n*<https://www.google.be|View on ADS-B>*",
					},
					{
						Type: "mrkdwn",
						Text: "*Callsign:*\nVIPER11",
					},
					{
						Type: "mrkdwn",
						Text: "*Type:*\nF16",
					},
				},
			},
			{
				Type: "divider",
			},
			{
				Type: "section",
				Fields: []Field{
					{
						Type: "mrkdwn",
						Text: "*Distance:*\n25 km",
					},
					{
						Type: "mrkdwn",
						Text: "*Tail number:*\nFA-117",
					},
				},
			},
		},
	}

	return slackMessage
}

func SendSlackMessage(aircraft []jetspotter.Aircraft) error {
	webHookURL := os.Getenv("SLACK_WEBHOOK_URL")

	data, err := json.Marshal(buildMessage(aircraft))
	if err != nil {
		return err
	}

	_, err = http.Post(webHookURL, "application/json",
		bytes.NewBuffer(data))
	if err != nil {
		return err
	}

	fmt.Println("A Slack notification has been sent!")
	return nil
}
