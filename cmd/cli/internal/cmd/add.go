package cmd

import (
	"fmt"

	"github.com/rs/zerolog/log"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/roloum/users/internal/user"
	"github.com/spf13/cobra"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Adds an user",
	//Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {

		email, _ := cmd.Flags().GetString("email")
		firstName, _ := cmd.Flags().GetString("first-name")
		lastName, _ := cmd.Flags().GetString("last-name")

		ctx := cmd.Context()
		log.Info().Msg("Executing the add command")

		cfg, ok := ctx.Value(ContextKey(CONFIG)).(Configuration)
		if !ok {
			return fmt.Errorf("Missing configuration")
		}

		dynamoDB, ok := ctx.Value(ContextKey(DYNAMO)).(*dynamodb.DynamoDB)
		if !ok {
			return fmt.Errorf("Missing DynamoDB connection")
		}

		nu := user.NewUser{
			Email:     email,
			FirstName: firstName,
			LastName:  lastName,
		}

		_, err := user.Create(ctx, dynamoDB, &nu, cfg.AWS.DynamoDB.Table.User)
		if err != nil {
			return err
		}

		log.Info().Msg("User created")
		return nil

	},
}

func init() {
	RootCmd.AddCommand(addCmd)

	var email, firstName, lastName string
	addCmd.Flags().StringVarP(&email, "email", "e", "", "Email (required)")
	addCmd.MarkFlagRequired("email")
	addCmd.Flags().StringVarP(&firstName, "first-name", "f", "", "First Name (required)")
	addCmd.MarkFlagRequired("first-name")
	addCmd.Flags().StringVarP(&lastName, "last-name", "l", "", "Last Name (required)")
	addCmd.MarkFlagRequired("last-name")
}
