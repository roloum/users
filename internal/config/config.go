package config

import (
	"log"

	"github.com/kelseyhightower/envconfig"
)

//Config holds the configuration for the application
type Config struct {
	AWS struct {
		DynamoDB struct {
			Table struct {
				User string `required:"true"`
			}
		}
	}
}

//Load Loads the configuration into the Config struct
func Load(log *log.Logger) (Config, error) {

	log.Println("Loading application configuration")
	var cfg Config

	err := envconfig.Process("users", &cfg)
	if err != nil {
		log.Printf("Error loading the configuration. %s\n", err.Error())
		return cfg, err
	}

	return cfg, nil

}
