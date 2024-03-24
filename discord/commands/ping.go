package commands

import (
	"GoBagouox/utils"
	"github.com/bwmarrin/discordgo"
)

func Ping(s *discordgo.Session, i *discordgo.InteractionCreate) {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Pong!",
		},
	})
	if err != nil {
		utils.Error("Cannot send pong message", err, 1)
	}

	message := &discordgo.MessageSend{
		Content: "Cliquez sur le bouton ci-dessous pour créer un btickets",
		Components: []discordgo.MessageComponent{
			&discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					&discordgo.Button{
						Label:    "Créer un btickets",
						CustomID: "ticket_create",
						Style:    discordgo.PrimaryButton,
					},
				},
			},
		},
	}

	_, err = s.ChannelMessageSendComplex(i.ChannelID, message)
	if err != nil {
		utils.Error("Failed to send message with button", err, 1)
	}
}
