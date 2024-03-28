package commands

import (
	"GoBagouox/utils"
	"github.com/bwmarrin/discordgo"
)

func RemoveChannel(s *discordgo.Session, i *discordgo.InteractionCreate) {
	embed := &discordgo.MessageEmbed{
		Title:       "Success",
		Description: "Removed",
		Color:       0x00ff00, // Vert
	}
	s.ChannelDelete(i.ChannelID)
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
