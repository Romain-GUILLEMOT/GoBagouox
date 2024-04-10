package tickets

import (
	"GoBagouox/database"
	"GoBagouox/database/models"
	"GoBagouox/discord/discordUtils"
	"GoBagouox/discord/service/tickets"
	"GoBagouox/utils"
	"github.com/bwmarrin/discordgo"
)

func Delete(s *discordgo.Session, i *discordgo.InteractionCreate) {
	channel := i.ChannelID
	db := database.GetDB()
	var ticket models.Ticket
	result := db.Where("channel_id = ?", channel).First(&ticket)
	if result.Error != nil {
		utils.Error("Try to close a unknow Ticket", result.Error, 1)
		discordUtils.DiscordError(s, i)
		return
	}
	//Get ticket
	err := tickets.Close(ticket, s)
	if err != nil {
		discordUtils.DiscordError(s, i)
	}
}
