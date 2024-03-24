package music

import (
	"GoBagouox/utils"
	"github.com/bwmarrin/discordgo"
)

func Play(s *discordgo.Session, i *discordgo.InteractionCreate) {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Pong!",
		},
	})
	if err != nil {
		utils.Error("Cannot send pong message", err, 1)
	}
}
