package notification

import (
	"jetspotter/internal/configuration"
	"jetspotter/internal/jetspotter"

	"github.com/bwmarrin/discordgo"
)

const (
	darkOrange  = 15755540
	lightOrange = 15829792
	darkYellow  = 15772952
	yellow      = 15055122
	lightGreen  = 10340365
	green       = 2278429
	greenBlue   = 1686636
	lightBlue   = 1292194
	darkBlue    = 2650083
	purple      = 10754265
	grey        = 3815994
)

// SendDiscordMessage sends a discord message containing metadata of a list of aircraft
func SendDiscordMessage(aircraft []jetspotter.Aircraft, config configuration.Config) error {
	// Discord has a limit of 10 embeds per message, so we need to split larger batches
	const maxEmbedsPerMessage = 10

	// If we have more than maxEmbedsPerMessage aircraft, split them into batches
	for i := 0; i < len(aircraft); i += maxEmbedsPerMessage {
		end := i + maxEmbedsPerMessage
		if end > len(aircraft) {
			end = len(aircraft)
		}

		batch := aircraft[i:end]
		message, err := buildDiscordMessage(batch, config)
		if err != nil {
			return err
		}

		notification := Notification{
			Message: message,
			Type:    Discord,
			URL:     config.DiscordWebHookURL,
		}

		err = SendMessage(notification)
		if err != nil {
			return err
		}
	}

	return nil
}

func getColorByAltitude(altitude int) int {
	switch {
	case altitude < 1000:
		return darkOrange
	case altitude >= 1000 && altitude < 2000:
		return lightOrange
	case altitude >= 2000 && altitude < 3000:
		return darkYellow
	case altitude >= 3000 && altitude < 5000:
		return yellow
	case altitude >= 5000 && altitude < 7000:
		return lightGreen
	case altitude >= 7000 && altitude < 10000:
		return green
	case altitude >= 10000 && altitude < 15000:
		return greenBlue
	case altitude >= 15000 && altitude < 20000:
		return lightBlue
	case altitude >= 20000 && altitude < 30000:
		return darkBlue
	case altitude >= 30000:
		return purple
	default:
		return grey
	}
}

func buildDiscordMessage(aircraft []jetspotter.Aircraft, config configuration.Config) (message discordgo.Message, err error) {
	message.Content = ":airplane: A jet has been spotted! :airplane:"
	var embeds []*discordgo.MessageEmbed
	for _, ac := range aircraft {
		embed := &discordgo.MessageEmbed{
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:   "Callsign",
					Value:  formatCallsign(ac, Markdown),
					Inline: true,
				},
				{
					Name:   "Registration",
					Value:  formatRegistration(ac, Markdown),
					Inline: true,
				},
				{
					Name:   "Country",
					Value:  ac.Country,
					Inline: true,
				},
				{
					Name:   "Speed",
					Value:  printSpeed(ac),
					Inline: true,
				},
				{
					Name:   "Altitude",
					Value:  printAltitude(ac),
					Inline: true,
				},
				{
					Name:   "Distance",
					Value:  printDistance(ac),
					Inline: true,
				},
				{
					Name:   "Bearing from location",
					Value:  printBearingFromLocation(ac),
					Inline: true,
				},
				{
					Name:   "Heading",
					Value:  printHeading(ac),
					Inline: true,
				},
				{
					Name:   "Bearing from aircraft",
					Value:  printBearingFromAircraft(ac),
					Inline: true,
				},
				{
					Name:   "Cloud coverage",
					Value:  printCloudCoverage(ac),
					Inline: true,
				},
				{
					Name:   "Inbound",
					Value:  getInboundStatus(ac),
					Inline: true,
				},
				{
					Name:   "Type",
					Value:  ac.Description,
					Inline: true,
				},
				{
					Name:   "Origin",
					Value:  printOriginName(ac),
					Inline: true,
				},
				{
					Name:   "Destination",
					Value:  printDestinationName(ac),
					Inline: true,
				},
				{
					Name:   "Airline",
					Value:  printAirlineName(ac),
					Inline: true,
				},
			},
		}

		if config.DiscordColorAltitude == "true" {
			embed.Color = getColorByAltitude(int(ac.Altitude))
		} else {
			embed.Color = darkBlue
		}

		imageURL := ac.ImageThumbnailURL
		if imageURL != "" {
			embed.Thumbnail = &discordgo.MessageEmbedThumbnail{
				URL: imageURL,
			}
		}

		embeds = append(embeds, embed)
	}

	message.Embeds = embeds
	return message, nil
}
