package commands

import (
	"GoBagouox/utils"
	"fmt"
	"github.com/bwmarrin/discordgo"
	probing "github.com/prometheus-community/pro-bing"
	"time"
)

type PingResult struct {
	Site    string
	Latency time.Duration
}

func getServerPing(site string, url string, results chan<- PingResult) {
	pinger, _ := probing.NewPinger(url)
	pinger.Count = 3
	err := pinger.Run()
	if err != nil {
		utils.Error("Can't ping a server.", err, 1)
	}
	stats := pinger.Statistics()

	results <- PingResult{
		Site:    site,
		Latency: stats.AvgRtt,
	}
}

func GetStatus(s *discordgo.Session, i *discordgo.InteractionCreate) {
	embed := &discordgo.MessageEmbed{
		Title:       "Success",
		Description: "Here is some stats about our system",
		Color:       0x00ff00, // Vert
	}
	apiLatency := int(s.HeartbeatLatency().Seconds() * 1000)
	Green := 0x00ff00
	Red := 0xff0000
	Yellow := 0xffff00
	fields := []*discordgo.MessageEmbedField{}

	sites := map[string]string{
		"API":    "api.bagou450.com",
		"CDN":    "cdn.bagou450.com",
		"DEMO":   "demo.bagou450.com",
		"HASTE":  "haste.bagou450.com",
		"UPLOAD": "UPLOAD.bagou450.com",
	}
	embed.Color = Green
	results := make(chan PingResult, len(sites))
	for site, url := range sites {
		go getServerPing(site, url, results)
	}
	for range sites {
		result := <-results

		fields = append(fields, &discordgo.MessageEmbedField{
			Name:   fmt.Sprintf("Status of %s", result.Site),
			Value:  fmt.Sprintf("%dms", result.Latency.Milliseconds()),
			Inline: true,
		})

		if embed.Color == Green {
			if result.Latency.Milliseconds() > 100 {
				embed.Color = Red
			} else if result.Latency.Milliseconds() > 50 {
				embed.Color = Yellow
			}
		} else if embed.Color == Yellow && result.Latency.Milliseconds() > 100 {
			embed.Color = Red
		}
	}
	fields = append(fields, &discordgo.MessageEmbedField{
		Name:   fmt.Sprintf("Status of Discord"),
		Value:  fmt.Sprintf("%dms", apiLatency),
		Inline: true,
	})
	embed.Fields = fields
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
		},
	})
	if err != nil {
		utils.Error("Cannot send pong message", err, 1)
	}
}
