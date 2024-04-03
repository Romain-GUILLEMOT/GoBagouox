package tickets

import (
	"GoBagouox/database"
	"GoBagouox/database/models"
	"GoBagouox/utils"
	"errors"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"gorm.io/gorm"
	"os"
)

func NewTicketMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	db := database.GetDB()
	var ticket models.Ticket
	result := db.First(&ticket, "channel_id = ?", m.ChannelID)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			utils.Debug("Channel "+m.ChannelID+" is not a ticket channel", 1)
			return
		}
		utils.Error("Can't save a ticket message.", result.Error, 1)
	}
	var user models.User
	result = db.First(&user, "discord_id = ?", m.Author.ID)
	if result.Error != nil {
		utils.Error("Can't find user with Discord ID "+m.Author.ID, result.Error, 1)
		return
	}

	ticketMessage := models.TicketMessage{
		Content:   m.Content,
		TicketID:  ticket.ID,
		MessageID: m.Message.ID,
		OwnerID:   user.ID,
	}

	result = db.Create(&ticketMessage)

	if result.Error != nil {
		utils.Error("Can't save a ticket message.", result.Error, 1)
		return
	}
	utils.Info("New ticket message saved.", 1)

	member, err := s.GuildMember(os.Getenv("DISCORD_GUILD"), m.Author.ID)
	if err != nil {
		utils.Error("Can't get user details.", result.Error, 1)
		return
	}
	isAdmin := false
	for _, role := range member.Roles {
		if role == os.Getenv("DISCORD_ADMIN_ROLE") {
			isAdmin = true
		}
	}
	data := map[string]string{
		"Link": fmt.Sprintf("https://discord.com/channels/%s/%s/%s", os.Getenv("DISCORD_GUILD"), m.ChannelID, m.ID),
	}
	if isAdmin && ticket.Status == "client_answer" {
		ticket.Status = "support_answer"
		db.Save(&ticket)
		go utils.SendEmail(user.Email, os.Getenv("APP_NAME")+" - New Message in Your Discord Ticket", "newTicketMessage", data)
	}

	if !isAdmin && ticket.Status == "support_answer" {
		ticket.Status = "client_answer"
		db.Save(&ticket)

		go utils.SendEmail(os.Getenv("OWNER_EMAIL"), os.Getenv("APP_NAME")+" - New Message in a Discord Ticket", "newTicketMessage", data)

	}

}
