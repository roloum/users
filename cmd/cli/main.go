package main

import (
	"context"
	"os"

	"github.com/roloum/users/cmd/cli/internal/cmd"
	uaws "github.com/roloum/users/internal/aws"
	"github.com/roloum/users/internal/config"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	zerolog.SetGlobalLevel(zerolog.DebugLevel)

	if err := run(); err != nil {
		log.Error().Msgf("Main: %s", err)
		os.Exit(1)
	}

}

func run() error {

	var cfg cmd.Configuration
	err := config.Load(&cfg)
	if err != nil {
		return err
	}
	ctx := context.WithValue(context.Background(), cmd.ContextKey(cmd.CONFIG), cfg)

	sess, err := uaws.GetSession(cfg.AWS.Region)
	if err != nil {
		return err
	}
	dynamo := uaws.GetDynamoDB(sess)
	ctx = context.WithValue(ctx, cmd.ContextKey(cmd.DYNAMO), dynamo)

	if err := cmd.RootCmd.ExecuteContext(ctx); err != nil {
		return err
	}

	return nil
}
