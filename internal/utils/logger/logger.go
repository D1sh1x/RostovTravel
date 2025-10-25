package logger

import (
	"os"

	"github.com/rs/zerolog"
)

func SetupLogger(env string) zerolog.Logger {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	var level zerolog.Level

	switch env {
	case "local":
		level = zerolog.DebugLevel
		zerolog.SetGlobalLevel(level)
		return zerolog.New(zerolog.NewConsoleWriter()).With().Timestamp().Logger()
	case "dev":
		level = zerolog.DebugLevel
	case "prod":
		level = zerolog.InfoLevel
	default:
		level = zerolog.InfoLevel
	}

	zerolog.SetGlobalLevel(level)
	return zerolog.New(os.Stdout).With().Timestamp().Logger()
}
