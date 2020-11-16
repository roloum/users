package cmd

import (
	"github.com/spf13/cobra"
)

//ContextKey ...
type ContextKey string

//LOG ...
const LOG = "log"

//AWSSESSION ...
const DYNAMO = "dynamo"

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "users",
	Short: "A CLI User Manager",
}
