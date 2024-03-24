package mtickets

import (
	"GoBagouox/utils"
	"encoding/json"
	"github.com/bwmarrin/discordgo"
	"log"
)

func Create(s *discordgo.Session, i *discordgo.InteractionCreate) {
	data := i.ModalSubmitData()
	components := data.Components
	jsonData, err := json.Marshal(components)
	if err != nil {
		log.Fatalf("Error while marshaling the Components: %v", err)
	}
	utils.Info(string(jsonData), 1)
	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Title:    "Title",
			CustomID: "your_custom_id",
			Content:  "Button clicked!",
		},
	})
	if err != nil {
		utils.Error("Cannot send a message.", err, 1)
	}
}
