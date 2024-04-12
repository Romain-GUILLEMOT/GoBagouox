package commands

import (
	"GoBagouox/utils"
	"github.com/bwmarrin/discordgo"
)

func GetYarn(s *discordgo.Session, i *discordgo.InteractionCreate) {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "You can retrieve your yarn logs with this command: \n ```js\nyarn build:production | curl -X POST -H \"Content-Type: text/plain\" --data-binary @- https://haste.bagou450.com/documents```",
		},
	})
	if err != nil {
		utils.Error("Cannot send pong message", err, 1)
	}

}
