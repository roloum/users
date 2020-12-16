package config

import (
	"os"

	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

//The init function sets the log level and format, since this is the file where
//all the configuration is loaded from the environment variables
func init() {

	if os.Getenv("USERS_LOG_PRETTY") == "pretty" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	//Set log level, default value which is Info
	var logLevel zerolog.Level
	switch os.Getenv("USERS_LOG_LEVEL") {
	case "fatal":
		logLevel = zerolog.FatalLevel
	case "error":
		logLevel = zerolog.ErrorLevel
	case "warn":
		logLevel = zerolog.WarnLevel
	case "info":
		logLevel = zerolog.InfoLevel
	case "debug":
		logLevel = zerolog.DebugLevel
	case "trace":
		logLevel = zerolog.TraceLevel
	default:
		logLevel = zerolog.InfoLevel
	}

	zerolog.SetGlobalLevel(logLevel)

	log.Debug().Msg("Log level set")
}

//Load Loads the configuration into the Config struct
func Load(cfg interface{}) error {

	log.Debug().Msg("Loading application configuration")

	err := envconfig.Process("users", cfg)
	if err != nil {
		log.Error().Msgf("Error loading the configuration. %s\n", err.Error())
		return err
	}

	return nil

}
