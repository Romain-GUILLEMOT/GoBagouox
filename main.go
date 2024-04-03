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
	"io"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"
)

func getLogFileWriter() io.Writer {
	// Set global logging level
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if os.Getenv("APP_DEBUG") == "true" {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		log.Debug().Msg("⚠️ " + utils.Yellow("Warning: Debug mode is enabled! This can reveal sensitive information on your web server ! Don't use it unless you know what you're doing.  ⚠️"))
	}

	// Set log time format
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	// Get logs folder
	logsFolder := os.Getenv("LOGS_FOLDER")
	if logsFolder == "" {
		logsFolder = "./logs"
	}

	// Create folder if not exist
	err := os.MkdirAll(logsFolder, os.ModePerm)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create logs dir")
	}

	// Get file name
	logFileName := time.Now().Format("2006-01-02") + ".log"
	logFilePath := filepath.Join(logsFolder, logFileName)

	// Create/Open the file
	logFile, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to open log file")
	}

	// Create a console writer with custom formatter for levels
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
	//Same for files
	fileWriter := zerolog.ConsoleWriter{Out: logFile, TimeFormat: zerolog.TimeFormatUnix}
	fileWriter.FormatLevel = func(i interface{}) string {
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
	// Return a multiwriter with log file and console writer
	return io.MultiWriter(fileWriter, consoleWriter)
}

func main() {
	//Load ENV
	err := godotenv.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("MAIN-001: Failed to load environment variables")
	}

	//Load logger
	writer := getLogFileWriter()
	log.Logger = log.Output(writer)

	//Test ENV
	iteration := os.Getenv("PBKDF2_ITERATIONS")
	salt := os.Getenv("SALT_SIZE")
	_, err = strconv.Atoi(iteration)
	if err != nil {
		utils.Fatal("Error during PBKDF2_ITERATIONS to int conversion.", err, 0)
	}
	_, err = strconv.Atoi(salt)
	if err != nil {
		utils.Fatal("Error during SALT_SIZE to int conversion.", err, 0)
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
