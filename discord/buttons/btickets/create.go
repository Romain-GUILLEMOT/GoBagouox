package btickets

import (
	"GoBagouox/database"
	"GoBagouox/database/models"
	"GoBagouox/discord/discordUtils"
	"GoBagouox/utils"
	"errors"
	"github.com/bwmarrin/discordgo"
	"gorm.io/gorm"
)

func Create(s *discordgo.Session, i *discordgo.InteractionCreate) {
	db := database.GetDB()
	form := []discordgo.MessageComponent{
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.TextInput{
					CustomID:    "subject",
					Label:       "Subject",
					Style:       discordgo.TextInputShort,
					Placeholder: "The ticket subject",
					Required:    true,
					MaxLength:   100,
					MinLength:   10,
				},
			},
		},
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.TextInput{
					CustomID:    "license",
					Label:       "License/Order ID",
					Placeholder: "Your license or order ID if you have one?",
					Style:       discordgo.TextInputShort,
					Required:    false,
					MaxLength:   100,
				},
			},
		},
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.TextInput{
					CustomID:    "logs",
					Label:       "Logs",
					Placeholder: "A link to your logs (Run /getlogs)?",
					Style:       discordgo.TextInputShort,
					Required:    false,
					MaxLength:   50,
					MinLength:   5,
				},
			},
		},
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.TextInput{
					CustomID:  "message",
					Label:     "Message",
					Style:     discordgo.TextInputParagraph,
					Required:  true,
					MinLength: 50,
					MaxLength: 2000,
				},
			},
		},
	}
	var user models.User
	result := db.First(&user, "discord_id = ?", i.Member.User.ID)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		form = append([]discordgo.MessageComponent{discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.TextInput{
					CustomID:    "emails",
					Label:       "Email",
					Placeholder: "emails@exemple.com",
					Style:       discordgo.TextInputShort,
					Required:    true,
					MinLength:   7,
				},
			},
		}}, form...)
	} else if result.Error != nil {
		utils.Error("Error during user check.", result.Error, 1)
		discordUtils.DiscordError(s, i)
		return
	}
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseModal,
		Data: &discordgo.InteractionResponseData{
			Title:      "Title",
			CustomID:   "ticket_create",
			Components: form,
		},
	})
	if err != nil {
		utils.Error("Cannot send a message.", err, 1)
	}
}
