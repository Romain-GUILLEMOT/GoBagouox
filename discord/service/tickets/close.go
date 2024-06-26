package tickets

import (
	"GoBagouox/database"
	"GoBagouox/database/models"
	"GoBagouox/utils"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"os"
	"strconv"
)

func Close(ticket models.Ticket, s *discordgo.Session) error {
	db := database.GetDB()
	//Set status to closed
	ticket.Status = "closed"
	result := db.Save(&ticket)
	if result.Error != nil {
		utils.Error("An unknown error occurred while saving of ticket modification.", result.Error, 1)
		return result.Error
	}
	//Send a email to the user
	user := ticket.Owner
	data := map[string]string{
		"Name":          ticket.Name,
		"Id":            strconv.Itoa(int(ticket.ID)),
		"Transcription": os.Getenv("MAINWEBSITE_URL") + "/account/ticket/discord/" + strconv.Itoa(int(ticket.ID)),
	}
	go utils.SendEmail(user.Email, os.Getenv("APP_NAME")+" - Ticket Closed", "ticket_closed", data)
	//Send private message to the user
	if user.DiscordID == "" {
		utils.Error("Discord ID is missing for the user", nil, 2)
	}
	userChannel, err := s.UserChannelCreate(user.DiscordID)
	utils.Debug("Create conversation with "+user.DiscordID, 1)
	if err != nil {
		utils.Error("An unknown error occured while create a conversation with the user", err, 1)
	}
	if userChannel != nil {
		userMention := "<@" + user.DiscordID + ">"
		message := fmt.Sprintf(
			"Hello %s,\n\n"+
				"Your ticket (%d - %s) has been closed.\n"+
				"You can find the ticket transcription on %s.\n\n"+
				"Best regards,\n"+
				"Bagouox\n"+
				"Bagou450 Team\n\n"+
				"PS: Notice that for checking your ticket transcription you need to have an account on our website with "+
				"the same email as the one provided during ticket creation.  :closed_book:",
			userMention, ticket.ID, ticket.Name,
			os.Getenv("MAINWEBSITE_URL")+"/account/ticket/discord/"+strconv.Itoa(int(ticket.ID)))

		_, err = s.ChannelMessageSend(userChannel.ID, message)
		if err != nil {
			utils.Error("An unknown error occurred while sending a private message to the user.", err, 1)
		}
	}

	//Remove channel
	_, err = s.ChannelDelete(ticket.ChannelId)
	if err != nil {
		utils.Error("An unknown error occurred while deleting the channel.", err, 1)
		return err
	}
	return nil
}
