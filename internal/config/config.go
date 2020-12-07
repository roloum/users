package config

import (
	"log"

	"github.com/kelseyhightower/envconfig"
)

//Load Loads the configuration into the Config struct
func Load(cfg interface{}, log *log.Logger) error {

	log.Println("Loading application configuration")

	err := envconfig.Process("users", cfg)
	if err != nil {
		log.Printf("Error loading the configuration. %s\n", err.Error())
		return err
	}

	return nil

}
