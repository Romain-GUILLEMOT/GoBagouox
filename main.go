package main

import (
	"GoBagouox/database"
	"GoBagouox/discord"
	"GoBagouox/utils"
	"GoBagouox/web"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"sync"
)

func main() {
	//Load logger
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	consoleWriter := zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: zerolog.TimeFormatUnix}
	consoleWriter.FormatLevel = func(i interface{}) string {
		switch i {
		case "debug":
			return fmt.Sprintf("| %-6s |", utils.Cyan("DEBUG"))
		case "info":
			return fmt.Sprintf("| %-6s |", utils.Green("INFO "))
		case "warn":
			return fmt.Sprintf("| %-6s |", utils.Yellow("WARN "))
		case "error":
			return fmt.Sprintf("| %-6s |", utils.Red("ERROR"))
		case "fatal":
			return fmt.Sprintf("| %-6s |", utils.BoldRed("FATAL"))
		default:
			return fmt.Sprintf("| %-6s |", i)
		}
	}

	log.Logger = log.Output(consoleWriter)

	//Load ENV
	err := godotenv.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("MAIN-001: Failed to load environment variables")
	}
	if os.Getenv("APP_DEBUG") == "true" {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		log.Debug().Msg("⚠️ " + utils.Yellow("Warning: Debug mode is enabled! This can reveal sensitive information on your web server ! Don't use it unless you know what you're doing.  ⚠️"))
	}
	var wg sync.WaitGroup
	// Init Database
	wg.Add(1)
	go func() {
		database.StartConnexion()
		wg.Done()
	}()

	// Init Discord server
	wg.Add(1)
	go func() {
		discord.StartBot()
		wg.Done()
	}()

	// Init Web server
	wg.Add(1)
	go func() {
		web.StartServer()
		wg.Done()
	}()

	wg.Wait()
}
