package notification

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"jetspotter/internal/jetspotter"
	"log"
	"net/http"
	"os"

	"github.com/PuerkitoBio/goquery"
	"github.com/jftuga/geodist"
)

type SlackMessage struct {
	Blocks []Block `json:"blocks"`
}

type Block struct {
	Type     string  `json:"type,omitempty"`
	Fields   []Field `json:"fields,omitempty"`
	Title    *Title  `json:"title,omitempty"`
	ImageURL string  `json:"image_url,omitempty"`
	AltText  string  `json:"alt_text,omitempty"`
}

type Title struct {
	Type  string `json:"type,omitempty"`
	Text  string `json:"text,omitempty"`
	Emoji bool   `json:"emoji,omitempty"`
}
type Field struct {
	Type string `json:"type,omitempty"`
	Text string `json:"text,omitempty"`
}

func buildMessage(aircraft []jetspotter.Aircraft, maxAmount int) (SlackMessage, error) {

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

		imageURL, err := getImageURL(fmt.Sprintf("https://www.planespotting.be/index.php?page=aircraft&registration=%s", ac.TailNumber))
		if imageURL != "" {
			blocks = append(blocks,
				Block{
					Type: "image",
					Title: &Title{
						Type:  "plain_text",
						Text:  fmt.Sprintf("%s - %s", ac.Desc, ac.TailNumber),
						Emoji: true,
					},
					ImageURL: imageURL,
					AltText:  fmt.Sprintf("%s with registration number %s", ac.Desc, ac.TailNumber),
				})
		}

		if err != nil {
			return SlackMessage{}, err
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

func getImageURL(URL string) (imageURL string, err error) {
	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		return "", err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}

	if res.StatusCode != 200 {
		return "", errors.New(fmt.Sprintf("Received status code %d", res.StatusCode))
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	doc.Find("meta[property='og:image']").Each(func(index int, element *goquery.Selection) {
		if content, exists := element.Attr("content"); exists {
			imageURL = content
		}
	})

	return imageURL, nil
}

func SendSlackMessage(aircraft []jetspotter.Aircraft, maxAmount int) error {
	webHookURL := os.Getenv("SLACK_WEBHOOK_URL")
	slackMessage, err := buildMessage(aircraft, maxAmount)
	if err != nil {
		return err
	}

	data, err := json.Marshal(slackMessage)
	if err != nil {
		return err
	}

	resp, err := http.Post(webHookURL, "application/json",
		bytes.NewBuffer(data))
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		fmt.Printf("%s\n", string(data))
		return errors.New(fmt.Sprintf("Received status code %v", resp.StatusCode))
	}

	fmt.Println("A Slack notification has been sent!")
	return nil
}
