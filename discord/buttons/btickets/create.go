package btickets

import (
	"GoBagouox/utils"
	"github.com/bwmarrin/discordgo"
)

func Create(s *discordgo.Session, i *discordgo.InteractionCreate) {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseModal,
		Data: &discordgo.InteractionResponseData{
			Title:    "Title",
			CustomID: "ticket_create",
			Content:  "Button clicked!",
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.TextInput{
							CustomID:    "opinion",
							Label:       "Subject?",
							Style:       discordgo.TextInputShort,
							Placeholder: "Don't be shy, share your opinion with us",
							Required:    true,
							MaxLength:   300,
							MinLength:   10,
						},
					},
				},
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.TextInput{
							CustomID:  "suggestions",
							Label:     "What would you suggest to improve them?",
							Style:     discordgo.TextInputParagraph,
							Required:  false,
							MaxLength: 2000,
						},
					},
				},
			},
		},
	})
	if err != nil {
		utils.Error("Cannot send a message.", err, 1)
	}
}
