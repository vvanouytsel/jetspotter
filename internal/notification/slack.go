package notification

import (
	"bytes"
	"encoding/json"
	"fmt"
	"jetspotter/internal/configuration"
	"jetspotter/internal/jetspotter"
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
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

func buildMessage(aircraft []jetspotter.AircraftOutput, config configuration.Config) (SlackMessage, error) {

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
		if i > config.MaxAircraftSlackMessage {
			break
		}

		blocks = append(blocks, Block{
			Type: "section",
			Fields: []Field{
				{
					Type: "mrkdwn",
					Text: fmt.Sprintf("*Callsign:* <https://globe.adsbexchange.com/?icao=%s | %s>", ac.ICAO, ac.Callsign),
				},
				{
					Type: "mrkdwn",
					Text: fmt.Sprintf("*Type:* %s", ac.Type),
				},
				{
					Type: "mrkdwn",
					Text: fmt.Sprintf("*Tail number:* %s", ac.TailNumber),
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
					Text: fmt.Sprintf("*Cloud coverage:* %d%%", ac.CloudCoverage),
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
						Text:  fmt.Sprintf("%s - %s", ac.Description, ac.TailNumber),
						Emoji: true,
					},
					ImageURL: imageURL,
					AltText:  fmt.Sprintf("%s with registration number %s", ac.Description, ac.TailNumber),
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
		fmt.Printf("Received status code %d for URL %s\n", res.StatusCode, URL)
		return "", nil
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

// SendSlackMessage sends a slack message containing metadata of a list of aircraft
func SendSlackMessage(aircraft []jetspotter.AircraftOutput, config configuration.Config) error {
	if len(aircraft) < 1 {
		return nil
	}

	slackMessage, err := buildMessage(aircraft, config)
	if err != nil {
		return err
	}

	data, err := json.Marshal(slackMessage)
	if err != nil {
		return err
	}

	resp, err := http.Post(config.SlackWebHookURL, "application/json",
		bytes.NewBuffer(data))
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		fmt.Printf("%s\n", string(data))
		return fmt.Errorf(fmt.Sprintf("Received status code %v", resp.StatusCode))
	}

	fmt.Println("A Slack notification has been sent!")
	return nil
}
