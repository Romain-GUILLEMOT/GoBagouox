package utils

import "github.com/rs/zerolog/log"

/*
Logtype:
0 = Web server
1 = Discord
2 = Database
*/
// Debug logging
func Debug(message string, logtype int) {
	if logtype == 1 {
		log.Debug().Msg(Blue("Discord: ") + message)
	} else if logtype == 2 {
		log.Debug().Msg(Brown("Database: ") + message)
	} else {
		log.Debug().Msg(Magenta("Web: ") + message)
	}
}

// Info logging
func Info(message string, logtype int) {
	if logtype == 1 {
		log.Info().Msg(Blue("Discord: ") + message)
	} else if logtype == 2 {
		log.Info().Msg(Brown("Database: ") + message)
	} else {
		log.Info().Msg(Magenta("Web: ") + message)
	}
}

// Warning logging
func Warning(message string, logtype int) {
	if logtype == 1 {
		log.Warn().Msg(Blue("Discord: ") + message)
	} else if logtype == 2 {
		log.Warn().Msg(Brown("Database: ") + message)
	} else {
		log.Warn().Msg(Magenta("Web: ") + message)
	}
}

// Error logging
func Error(message string, err error, logtype int) {
	if logtype == 1 {
		log.Error().Err(err).Msg(Blue("Discord: ") + message)
	} else if logtype == 2 {
		log.Error().Err(err).Msg(Brown("Database: ") + message)
	} else {
		log.Error().Err(err).Msg(Magenta("Web: ") + message)
	}
}

// Fatal logging
func Fatal(message string, err error, logtype int) {
	if logtype == 1 {
		log.Fatal().Err(err).Msg(Blue("Discord: ") + message)
	} else if logtype == 2 {
		log.Fatal().Err(err).Msg(Brown("Database: ") + message)
	} else {
		log.Fatal().Err(err).Msg(Magenta("Web: ") + message)
	}
}
