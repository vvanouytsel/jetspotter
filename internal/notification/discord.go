package notification

import (
	"fmt"
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
func SendDiscordMessage(aircraft []jetspotter.AircraftOutput, config configuration.Config) error {
	message, err := buildDiscordMessage(aircraft, config)
	notification := Notification{
		Message:    message,
		Type:       Discord,
		WebHookURL: config.DiscordWebHookURL,
	}

	err = SendMessage(aircraft, notification)
	if err != nil {
		return err
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

func buildDiscordMessage(aircraft []jetspotter.AircraftOutput, config configuration.Config) (message discordgo.Message, err error) {
	message.Content = ":airplane: A jet has been spotted! :airplane:"
	var embeds []*discordgo.MessageEmbed
	for _, ac := range aircraft {
		embed := &discordgo.MessageEmbed{
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:   "Callsign",
					Value:  fmt.Sprintf("[%s](%s)", ac.Callsign, ac.TrackerURL),
					Inline: true,
				},
				{
					Name:   "Tail number",
					Value:  ac.TailNumber,
					Inline: true,
				},
				{
					Name:   "Speed",
					Value:  fmt.Sprintf("%dkn | %dkm/h", ac.Speed, jetspotter.ConvertKnotsToKilometersPerHour(ac.Speed)),
					Inline: true,
				},
				{
					Name:   "Altitude",
					Value:  fmt.Sprintf("%vft | %dm", ac.Altitude, jetspotter.ConvertFeetToMeters(ac.Altitude)),
					Inline: true,
				},
				{
					Name:   "Distance",
					Value:  fmt.Sprintf("%dkm", ac.Distance),
					Inline: true,
				},
				{
					Name:   "Bearing from location",
					Value:  fmt.Sprintf("%.0f°", ac.BearingFromLocation),
					Inline: true,
				},
				{
					Name:   "Heading",
					Value:  fmt.Sprintf("%.0f°", ac.Heading),
					Inline: true,
				},
				{
					Name:   "Bearing from aircraft",
					Value:  fmt.Sprintf("%.0f°", ac.BearingFromAircraft),
					Inline: true,
				},
				{
					Name:   "Cloud coverage",
					Value:  fmt.Sprintf("%d%%", ac.CloudCoverage),
					Inline: true,
				},
				{
					Name:   "Type",
					Value:  ac.Description,
					Inline: true,
				},
			},
		}

		if config.DiscordColorAltitude == "true" {
			embed.Color = getColorByAltitude(int(ac.Altitude))
		} else {
			embed.Color = darkBlue
		}

		imageURL := ac.PlaneSpotterURL
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
