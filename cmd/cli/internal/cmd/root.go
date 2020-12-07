package cmd

import (
	"github.com/spf13/cobra"
)

//ContextKey ...
type ContextKey string

//LOG Application Log
const LOG = "log"

//CONFIG Application configuration struct
const CONFIG = "config"

//DYNAMO AWS DynamoDB session
const DYNAMO = "dynamo"

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "users",
	Short: "A CLI User Manager",
}

//Configuration stores the configuration for the cli commads
type Configuration struct {
	AWS struct {
		DynamoDB struct {
			Table struct {
				User string `required:"true"`
			}
		}
	}
}
