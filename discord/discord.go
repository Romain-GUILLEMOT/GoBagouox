package discord

import (
	"GoBagouox/discord/buttons/btickets"
	"GoBagouox/discord/commands"
	"GoBagouox/discord/message/tickets"
	"GoBagouox/discord/modals/mtickets"
	"GoBagouox/utils"
	"github.com/bwmarrin/discordgo"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type CommandFunc func(s *discordgo.Session, i *discordgo.InteractionCreate)
type NewMessageFunc func(s *discordgo.Session, m *discordgo.MessageCreate)

type Command struct {
	Name        string
	Description string
	Handler     CommandFunc
	Options     []*discordgo.ApplicationCommandOption
}
type NewMessage struct {
	Name    string
	Handler NewMessageFunc
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
		"removechannel": {
			Name:        "removechannel",
			Description: "removechannel",
			Handler:     commands.RemoveChannel,
		},
	}
}
func getButtons() map[string]Command {
	return map[string]Command{
		"ticket_create": {
			Name:        "ticket_create",
			Description: "Create a ticket!",
			Handler:     btickets.Create,
		},
		"ticket_close": {
			Name:        "ticket_close",
			Description: "Close a ticket!",
			Handler:     btickets.Delete,
		},
	}
}
func getModals() map[string]Command {
	return map[string]Command{
		"ticket_create": {
			Name:        "ticket_create",
			Description: "Create a ticket!",
			Handler:     mtickets.Create,
		},
	}
}
func getNewMessage() map[string]NewMessage {
	return map[string]NewMessage{
		"ticketsMessage": {
			Name:    "ticketsMessage",
			Handler: tickets.NewTicketMessage,
		},
		"ticketsAttachments": {
			Name:    "ticketsAttachments",
			Handler: tickets.NewTicketAttachment,
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
	registerButtons()
	registerModals()
	dg.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if m.Author.ID == s.State.User.ID {
			return
		}
		utils.Debug("New message. Run all functions.", 1)

		for _, newmesshandler := range getNewMessage() {
			utils.Debug("Run: "+newmesshandler.Name, 1)
			newmesshandler.Handler(s, m)
		}
	})
	dg.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {

		if i.Type == discordgo.InteractionApplicationCommand {
			commandList := getCommands()
			command, commandExists := commandList[i.ApplicationCommandData().Name]
			if commandExists {

				utils.Info("Execution of the command "+utils.Bold(utils.Blue(command.Name)), 1)
				startTime := time.Now()
				command.Handler(s, i)
				executionTime := time.Since(startTime)
				utils.Info("End of execution of the command "+utils.Bold(utils.Blue(command.Name))+" in "+executionTime.String(), 1)

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
		if i.Type == discordgo.InteractionMessageComponent {
			buttonList := getButtons()
			button, buttonExists := buttonList[i.MessageComponentData().CustomID]
			if buttonExists {
				utils.Info("Execution of the button "+utils.Bold(utils.Blue(button.Name)), 1)
				startTime := time.Now()
				button.Handler(s, i)
				executionTime := time.Since(startTime)
				utils.Info("End of execution of the button "+utils.Bold(utils.Blue(button.Name))+" in "+executionTime.String(), 1)

			} else {
				embed := &discordgo.MessageEmbed{
					Title:       "Error",
					Description: "This button (" + i.MessageComponentData().CustomID + ") was not found in our database. Please try again.",
					Color:       0xff0000, // Rouge
				}
				err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Embeds: []*discordgo.MessageEmbed{embed},
					},
				})
				if err != nil {
					utils.Error("Cannot send a message.", err, 1)
				}
			}
		}
		if i.Type == discordgo.InteractionModalSubmit {
			modalsList := getModals()
			modal, modalExists := modalsList[i.ModalSubmitData().CustomID]
			if modalExists {
				utils.Info("Execution of the modal "+utils.Bold(utils.Blue(modal.Name)), 1)
				startTime := time.Now()
				modal.Handler(s, i)
				executionTime := time.Since(startTime)
				utils.Info("End of execution of the modal "+utils.Bold(utils.Blue(modal.Name))+" in "+executionTime.String(), 1)

			} else {
				embed := &discordgo.MessageEmbed{
					Title:       "Error",
					Description: "This modal (" + i.MessageComponentData().CustomID + ") was not found in our database. Please try again.",
					Color:       0xff0000, // Rouge
				}
				err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Embeds: []*discordgo.MessageEmbed{embed},
					},
				})
				if err != nil {
					utils.Error("Cannot send a message.", err, 1)
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
			Options:     command.Options,
		}

		_, err := session.ApplicationCommandCreate(session.State.User.ID, guildID, discordCommand)
		if err != nil {
			utils.Error("ERR-003: Failed to create slash command: "+utils.Bold(discordCommand.Name), err, 1)
		} else {
			utils.Debug("Loaded slash command: "+utils.Bold(utils.Blue(discordCommand.Name)), 1)
		}
	}
}
func registerButtons() {
	buttonsList := getButtons()

	for _, button := range buttonsList {
		utils.Debug("Loaded button : "+utils.Bold(utils.Blue(button.Name)), 1)

	}
}
func registerModals() {
	modalsList := getModals()

	for _, modal := range modalsList {
		utils.Debug("Loaded modal : "+utils.Bold(utils.Blue(modal.Name)), 1)

	}
}
