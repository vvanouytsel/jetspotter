package notification

import (
	"fmt"
	"jetspotter/internal/configuration"
	"jetspotter/internal/jetspotter"
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/gotify/go-api-client/v2/models"
)

// OverheadContext bundles the information the per-platform overhead message
// builders need in addition to the aircraft snapshot. Keeping it in one
// struct makes adding new platforms and rendering consistent.
type OverheadContext struct {
	Candidate jetspotter.OverheadCandidate
}

// headline produces the natural-language "look up" headline shared by all
// platforms:
//   "In 3 min, look up — you'll see a B748 heading south."
// When the minutes value is 0 it reads "look up now".
func headline(c jetspotter.OverheadCandidate) string {
	aircraftKind := c.Aircraft.Type
	if aircraftKind == "" {
		aircraftKind = c.Aircraft.Description
	}
	if aircraftKind == "" {
		aircraftKind = "aircraft"
	}

	if c.MinutesUntilOverhead <= 0 {
		return fmt.Sprintf("Look up now! A %s is passing overhead heading %s.", aircraftKind, c.CompassWord)
	}
	return fmt.Sprintf("In %d min, look up — you'll see a %s heading %s.", c.MinutesUntilOverhead, aircraftKind, c.CompassWord)
}

// SendOverheadTerminalMessage prints the overhead prediction to the terminal.
func SendOverheadTerminalMessage(c jetspotter.OverheadCandidate, config configuration.Config) {
	log.Println("👀 Overhead prediction!")
	fmt.Println(headline(c))
	fmt.Println(FormatAircraft(c.Aircraft, config))
}

// SendOverheadSlackMessage sends the overhead prediction to Slack.
func SendOverheadSlackMessage(c jetspotter.OverheadCandidate, config configuration.Config) error {
	blocks := []Block{
		{
			Type: "section",
			Fields: []Field{
				{Type: "mrkdwn", Text: fmt.Sprintf(":eyes: %s", headline(c))},
			},
		},
		{
			Type: "section",
			Fields: []Field{
				{Type: "mrkdwn", Text: fmt.Sprintf("*Callsign:* <%s|%s>", c.Aircraft.TrackerURL, c.Aircraft.Callsign)},
				{Type: "mrkdwn", Text: fmt.Sprintf("*Type:* %s", c.Aircraft.Type)},
				{Type: "mrkdwn", Text: fmt.Sprintf("*Altitude:* %s", printAltitude(c.Aircraft))},
				{Type: "mrkdwn", Text: fmt.Sprintf("*Speed:* %s", printSpeed(c.Aircraft))},
			},
		},
	}

	message := SlackMessage{Blocks: blocks}
	notification := Notification{
		Message: message,
		Type:    Slack,
		URL:     config.SlackWebHookURL,
	}
	return SendMessage(notification)
}

// SendOverheadDiscordMessage sends the overhead prediction to Discord.
func SendOverheadDiscordMessage(c jetspotter.OverheadCandidate, config configuration.Config) error {
	aircraftKind := c.Aircraft.Type
	if aircraftKind == "" {
		aircraftKind = c.Aircraft.Description
	}

	embed := &discordgo.MessageEmbed{
		Title:       "👀 Overhead Prediction",
		Description: headline(c),
		Color:       darkBlue,
		Fields: []*discordgo.MessageEmbedField{
			{Name: "Callsign", Value: formatCallsign(c.Aircraft, Markdown), Inline: true},
			{Name: "Type", Value: aircraftKind, Inline: true},
			{Name: "Heading", Value: fmt.Sprintf("%s (%.0f°)", c.CompassWord, c.Aircraft.Heading), Inline: true},
			{Name: "Altitude", Value: printAltitude(c.Aircraft), Inline: true},
			{Name: "Speed", Value: printSpeed(c.Aircraft), Inline: true},
			{Name: "CPA Distance", Value: fmt.Sprintf("%dkm", c.CPADistanceKm), Inline: true},
		},
	}

	if c.Aircraft.ImageThumbnailURL != "" {
		embed.Thumbnail = &discordgo.MessageEmbedThumbnail{URL: c.Aircraft.ImageThumbnailURL}
	}

	message := discordgo.Message{
		Content: ":eyes: A plane is about to fly overhead!",
		Embeds:  []*discordgo.MessageEmbed{embed},
	}

	notification := Notification{
		Message: message,
		Type:    Discord,
		URL:     config.DiscordWebHookURL,
	}
	return SendMessage(notification)
}

// SendOverheadGotifyMessage sends the overhead prediction to Gotify.
func SendOverheadGotifyMessage(c jetspotter.OverheadCandidate, config configuration.Config) error {
	message := models.MessageExternal{
		Title: "Overhead prediction!",
		Extras: map[string]interface{}{
			"client::display": map[string]interface{}{
				"contentType": "text/markdown",
			},
		},
	}
	message.Message = fmt.Sprintf("**%s**\n\n", headline(c))
	message.Message += fmt.Sprintf("**Callsign**: %s\n\n", formatCallsign(c.Aircraft, Markdown))
	message.Message += fmt.Sprintf("**Type**: %s\n\n", c.Aircraft.Type)
	message.Message += fmt.Sprintf("**Altitude**: %s\n\n", printAltitude(c.Aircraft))
	message.Message += fmt.Sprintf("**Speed:** %s\n\n", printSpeed(c.Aircraft))
	message.Message += fmt.Sprintf("**CPA distance:** %dkm\n\n", c.CPADistanceKm)

	notification := Notification{
		Message: message,
		Type:    Gotify,
		URL:     fmt.Sprintf("%s/message?token=%s", config.GotifyURL, config.GotifyToken),
	}
	return SendMessage(notification)
}

// SendOverheadNtfyMessage sends the overhead prediction to ntfy.
func SendOverheadNtfyMessage(c jetspotter.OverheadCandidate, config configuration.Config) error {
	message := NtfyNotification{
		Title:    "Overhead prediction!",
		Topic:    config.NtfyTopic,
		Tags:     []string{"eyes", "airplane"},
		Markdown: true,
		Actions: []NtfyAction{
			AddNtfyAction("Track Aircraft", c.Aircraft.TrackerURL),
		},
	}
	message.Message = fmt.Sprintf("%s\n\n", headline(c))
	message.Message += fmt.Sprintf("Callsign: %s\n", formatCallsign(c.Aircraft, Markdown))
	message.Message += fmt.Sprintf("Type: %s\n", c.Aircraft.Type)
	message.Message += fmt.Sprintf("Altitude: %s\n", printAltitude(c.Aircraft))
	message.Message += fmt.Sprintf("Speed: %s\n", printSpeed(c.Aircraft))
	message.Message += fmt.Sprintf("CPA distance: %dkm\n", c.CPADistanceKm)

	notification := Notification{
		Message: message,
		Type:    Ntfy,
		URL:     config.NtfyServer,
		Token:   config.NtfyToken,
	}
	return SendMessage(notification)
}
