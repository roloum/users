package cmd

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/roloum/users/internal/user"
)

var activateCmd = &cobra.Command{
	Use:   "activate",
	Short: "Activates an user",
	RunE: func(cmd *cobra.Command, args []string) error {

		email, _ := cmd.Flags().GetString("email")

		ctx := cmd.Context()
		log.Info().Msg("Executing the activate command")

		cfg, ok := ctx.Value(ContextKey(CONFIG)).(Configuration)
		if !ok {
			return fmt.Errorf("Missing configuration")
		}

		dynamoDB, ok := ctx.Value(ContextKey(DYNAMO)).(*dynamodb.DynamoDB)
		if !ok {
			return fmt.Errorf("Missing DynamoDB connection")
		}

		u := &user.User{
			Email: email,
		}

		if err := u.Activate(ctx, dynamoDB, cfg.AWS.DynamoDB.Table.User); err != nil {
			log.Error().Msg(err.Error())
			return err
		}

		log.Info().Msg("User activated")

		return nil
	},
}

func init() {
	RootCmd.AddCommand(activateCmd)

	var email string
	activateCmd.Flags().StringVarP(&email, "email", "e", "", "Email (required)")
	activateCmd.MarkFlagRequired("email")
}
