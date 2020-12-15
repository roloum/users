package config

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog/log"
)

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
