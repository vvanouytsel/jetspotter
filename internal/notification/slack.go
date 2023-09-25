package notification

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"jetspotter/internal/jetspotter"
	"net/http"
	"os"

	"github.com/jftuga/geodist"
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

func buildMessage(aircraft []jetspotter.Aircraft, maxAmount int) SlackMessage {

	var blocks []Block
	blocks = append(blocks, Block{
		Type: "section",
		Fields: []Field{
			{
				Type: "mrkdwn",
				Text: ":jet: A jet has been spotted! :jet:",
			},
		},
	})

	for i, ac := range aircraft {
		if i > maxAmount {
			break
		}

		blocks = append(blocks, Block{
			Type: "section",
			Fields: []Field{
				{
					Type: "mrkdwn",
					Text: fmt.Sprintf("*Callsign:* <https://globe.adsbexchange.com/?icao=%s| %s>", ac.ICAO, ac.Callsign),
				},
				{
					Type: "mrkdwn",
					Text: fmt.Sprintf("*Type:* %s", ac.PlaneType),
				},
				{
					Type: "mrkdwn",
					Text: fmt.Sprintf("*Tail number:* %s", ac.TailNumber),
				},
				{
					Type: "mrkdwn",
					Text: fmt.Sprintf("*Speed:* %dkn", int(ac.GS)),
				},
				{
					Type: "mrkdwn",
					Text: fmt.Sprintf("*Altitude:* %vft", ac.AltBaro),
				},
				{
					Type: "mrkdwn",
					Text: fmt.Sprintf("*Distance:* %vkm", jetspotter.CalculateDistanceToBullseye(geodist.Coord{
						Lat: ac.Lat,
						Lon: ac.Lon,
					})),
				},
			},
		})

		blocks = append(blocks, Block{
			Type: "divider",
		})
	}
	//				"text": "You have a new request:\n*<fakeLink.toEmployeeProfile.com|Fred Enriquez - New device request>*"

	slackMessage := SlackMessage{Blocks: blocks}
	return slackMessage
}

func SendSlackMessage(aircraft []jetspotter.Aircraft, maxAmount int) error {
	webHookURL := os.Getenv("SLACK_WEBHOOK_URL")

	data, err := json.Marshal(buildMessage(aircraft, maxAmount))
	if err != nil {
		return err
	}

	resp, err := http.Post(webHookURL, "application/json",
		bytes.NewBuffer(data))
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		fmt.Printf("%s", string(data))
		return errors.New(fmt.Sprintf("Received status code %v", resp.StatusCode))
	}

	fmt.Println("A Slack notification has been sent!")
	return nil
}
