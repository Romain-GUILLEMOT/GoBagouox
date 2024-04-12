package tickets

import (
	"GoBagouox/database"
	"GoBagouox/database/models"
	"GoBagouox/utils"
	"github.com/bwmarrin/discordgo"
	"strings"
)

func SendMessageToTicket(ticket models.Ticket, s *discordgo.Session, content string, user models.User) error {
	db := database.GetDB()
	name := user.Username
	if strings.ToLower(name) == "bagou450" || strings.ToLower(name) == "bagou450second" {
		name = "Romain GUILLEMOT"
	}
	message := "New message sent from Website by **" + name + "**:\n" + content

	discordMessage, err := s.ChannelMessageSend(ticket.ChannelId, message)
	if err != nil {
		utils.Error("An error ocurred during sending message to ticket process", err, 1)
		return err
	}
	ticketMessage := models.TicketMessage{
		Content:   content,
		MessageID: discordMessage.ID,
		Ticket:    ticket,
		Owner:     user,
	}
	db.Save(&ticketMessage)
	return nil
}
