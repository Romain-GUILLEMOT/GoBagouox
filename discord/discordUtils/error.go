package discordUtils

import (
	"GoBagouox/utils"
	"github.com/bwmarrin/discordgo"
)

func DiscordError(s *discordgo.Session, i *discordgo.InteractionCreate) {
	embed := &discordgo.MessageEmbed{
		Title:       "ERROR",
		Description: "An unknown error occurred. Please contact a member of our support team (contact@bagou450.com).",
		Color:       0xFF0000,
	}
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
		},
	})
	if err != nil {
		utils.Error("Cannot send Error message", err, 1)
	}
}
