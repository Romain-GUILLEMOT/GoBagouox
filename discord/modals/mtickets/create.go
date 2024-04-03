package mtickets

import (
	"GoBagouox/database"
	"GoBagouox/database/models"
	"GoBagouox/discord/discordUtils"
	"GoBagouox/utils"
	"encoding/json"
	"errors"
	"github.com/bwmarrin/discordgo"
	"gorm.io/gorm"
	"os"
	"regexp"
	"strconv"
)

var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

type Component struct {
	CustomID string `json:"custom_id"`
	Value    string `json:"value"`
}
type ComponentData struct {
	Components []Component `json:"components"`
}

func Create(s *discordgo.Session, i *discordgo.InteractionCreate) {
	data := i.ModalSubmitData()
	components := data.Components
	var user models.User
	license := "No license provided"
	logs := "No logs provided"
	subject := ""
	message := ""
	db := database.GetDB()
	result := db.Where("discord_id = ?", i.Member.User.ID).First(&user)
	if result.Error != nil {
		if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			utils.Error("Unknow database error. ", result.Error, 2)
			discordUtils.DiscordError(s, i)
			return
		}
	}
	for _, item := range components {
		itemJson, err := item.MarshalJSON()
		if err != nil {
			return
		}
		var data ComponentData
		err = json.Unmarshal(itemJson, &data)
		firstComponent := data.Components[0]
		switch firstComponent.CustomID {
		case "subject":
			subject = firstComponent.Value
		case "message":
			message = firstComponent.Value
		case "license":
			if firstComponent.Value != "" {
				license = firstComponent.Value
			}
		case "logs":
			if firstComponent.Value != "" {
				logs = firstComponent.Value
			}
		}
		if firstComponent.CustomID == "emails" && firstComponent.Value != "" {
			if !emailRegex.MatchString(firstComponent.Value) {
				embed := &discordgo.MessageEmbed{
					Title:       "ERROR",
					Description: "An invalid emails was provided.",
					Color:       0xFF0000,
				}
				err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Embeds: []*discordgo.MessageEmbed{embed},
						Flags:  discordgo.MessageFlagsEphemeral,
					},
				})
				if err != nil {
					utils.Error("Cannot send a message", err, 1)

				}
				return
			}
			user = models.User{
				Model:     gorm.Model{},
				DiscordID: i.Member.User.ID,
				Username:  i.Member.User.Username,
				Avatar:    i.Member.User.AvatarURL("512"),
				Email:     firstComponent.Value,
			}
			db.Create(&user)

			utils.Info("New user "+firstComponent.Value+" was created!", 1)
		}
	}
	var ticket models.Ticket
	db.Order("id desc").First(&ticket)

	channel, err := s.GuildChannelCreate(os.Getenv("DISCORD_GUILD"), "ticket-"+strconv.Itoa(int(ticket.ID+1)), discordgo.ChannelTypeGuildText)
	//channel, err := s.GuildChannelCreate(os.Getenv("DISCORD_GUILD"), "ticket-"+strconv.Itoa(int(1)), discordgo.ChannelTypeGuildText)
	if err != nil {
		utils.Error("Cannot create ticket channel.", err, 1)
		discordUtils.DiscordError(s, i)
		return
	}
	perms := int64(0x0000400000000000 | discordgo.PermissionReadMessageHistory | discordgo.PermissionViewChannel | discordgo.PermissionUseExternalEmojis | discordgo.PermissionSendMessages | discordgo.PermissionEmbedLinks | discordgo.PermissionAttachFiles | discordgo.PermissionAddReactions)

	err = s.ChannelPermissionSet(
		channel.ID,
		os.Getenv("DISCORD_EVERYONE_ROLE"),
		discordgo.PermissionOverwriteTypeRole,
		0,
		perms,
	)
	if err != nil {
		utils.Error("Cannot edit permission of ticket channel.", err, 1)
		discordUtils.DiscordError(s, i)
		return
	}
	err = s.ChannelPermissionSet(
		channel.ID,
		user.DiscordID,                          // l'ID du rôle auquel vous voulez affecter ces permissions
		discordgo.PermissionOverwriteTypeMember, // le type de l'ID ci-dessus. Cela peut être "member" pour un utilisateur spécifique ou "role" pour un rôle entier
		perms,                                   // les permissions que vous voulez autoriser. Voir : https://discord.com/developers/docs/topics/permissions#permissions-bitwise-permission-flags
		0,                                       // les permissions que vous voulez refuser. Voir : https://discord.com/developers/docs/topics/permissions#permissions-bitwise-permission-flags
	)
	if err != nil {
		utils.Error("Cannot edit permission of ticket channel.", err, 1)
		discordUtils.DiscordError(s, i)
		return
	}
	ticket = models.Ticket{
		Model:     gorm.Model{},
		Name:      subject,
		Status:    "client_answer",
		License:   license,
		Logs:      logs,
		ChannelId: channel.ID,
		Owner:     user,
	}
	result = db.Create(&ticket)
	if result.Error != nil {
		utils.Error("Cannot create ticket.", result.Error, 1)
		discordUtils.DiscordError(s, i)
		return
	}
	deleteComponent := []discordgo.MessageComponent{
		&discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				&discordgo.Button{
					Label:    "Close ticket",
					CustomID: "ticket_close",
					Style:    discordgo.DangerButton,
				},
			},
		},
	}
	embed := &discordgo.MessageEmbed{
		Title:       subject,
		Description: "This ticket was made by <@" + user.DiscordID + "> (" + user.Email + ")\n With these informations:",
		Color:       0x00ff00, // Vert
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Logs",
				Value:  logs,
				Inline: false,
			},
			{
				Name:   "License",
				Value:  license,
				Inline: false,
			},
		},
	}
	ticketMessage := models.TicketMessage{
		Model:   gorm.Model{},
		Content: message,
		Ticket:  ticket,
		Owner:   user,
	}
	db.Create(&ticketMessage)

	_, err = s.ChannelMessageSendComplex(channel.ID, &discordgo.MessageSend{
		Embed:      embed,
		Components: deleteComponent,
	})
	embed = &discordgo.MessageEmbed{
		Title:       "Success",
		Description: "Your ticket #ticket-" + strconv.Itoa(int(ticket.ID)) + " has been created",
		Color:       0x00ff00,
	}
	if err != nil {
		utils.Error("Cannot send a message to the ticket channel.", err, 1)
		embed = &discordgo.MessageEmbed{
			Title:       "Warning",
			Description: "The ticket was created but we can't send any message on it. \nPlease write your message in the ticket again.",
			Color:       0xFFFF00,
		}
	}
	_, err = s.ChannelMessageSend(channel.ID, "Provided message: \n"+message)
	if err != nil {
		utils.Error("Cannot send a message to the ticket channel.", err, 1)
		embed = &discordgo.MessageEmbed{
			Title:       "Warning",
			Description: "The ticket was created but we can't send any message on it. \nPlease write your message in the ticket again.",
			Color:       0xFFFF00,
		}
	}

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
		},
	})
	if err != nil {
		utils.Error("Cannot send a message", err, 1)
	}

}
