package discord

import (
	"GoBagouox/discord/commands"
	"GoBagouox/discord/commands/music"
	"GoBagouox/utils"
	"github.com/bwmarrin/discordgo"
	"os"
	"os/signal"
	"syscall"
)

type CommandFunc func(s *discordgo.Session, i *discordgo.InteractionCreate)

type Command struct {
	Name        string
	Description string
	Handler     CommandFunc
}

func getCommands() map[string]Command {
	return map[string]Command{
		"ping": {
			Name:        "ping",
			Description: "Answer Pong!",
			Handler:     commands.Ping,
		},
		"status": {
			Name:        "status",
			Description: "Return Bagou450 server status.",
			Handler:     commands.GetStatus,
		},
		"play": {
			Name:        "play",
			Description: "Play a youtube music.",
			Handler:     music.Play,
		},
	}
}
func StartBot() {

	dg, err := discordgo.New("Bot " + os.Getenv("DISCORD_TOKEN"))
	if err != nil {
		utils.Fatal("ERR-001: Failed to create Discord session", err, 1)
	}

	err = dg.Open()
	if err != nil {
		utils.Fatal("ERR-002: Failed to open websocket connection to Discord", err, 1)
	}

	registerSlashCommands(dg, os.Getenv("DISCORD_GUILD"))
	dg.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if i.Type == discordgo.InteractionApplicationCommand {
			commandList := getCommands()
			command, commandExists := commandList[i.ApplicationCommandData().Name]
			if commandExists {
				command.Handler(s, i)
			} else {
				embed := &discordgo.MessageEmbed{
					Title:       "Error",
					Description: "This command was not found in our database. Please try again.",
					Color:       0xff0000, // Rouge
				}
				err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Embeds: []*discordgo.MessageEmbed{embed},
					},
				})
				if err != nil {
					utils.Error("Cannot send command not found message", err, 1)
				}
			}
		}
	})

	utils.Info("Discord bot is now running.", 1)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	utils.Info("Discord bot has now stopped.", 1)
	err = dg.Close()
	if err != nil {
		utils.Error("Cannot stop Discord bot.", err, 1)

	}

}

func registerSlashCommands(session *discordgo.Session, guildID string) {
	commandsList := getCommands()

	for _, command := range commandsList {
		discordCommand := &discordgo.ApplicationCommand{
			Name:        command.Name,
			Description: command.Description,
		}

		_, err := session.ApplicationCommandCreate(session.State.User.ID, guildID, discordCommand)
		if err != nil {
			utils.Error("ERR-003: Failed to create slash command: "+utils.Bold(discordCommand.Name), err, 1)
		} else {
			utils.Debug("Loaded slash command: "+utils.Bold(utils.Blue(discordCommand.Name)), 1)
		}
	}
}
