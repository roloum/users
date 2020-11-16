package main

import (
	"context"
	"log"

	"os"

	"github.com/roloum/users/cmd/cli/internal/cmd"
	uaws "github.com/roloum/users/internal/aws"
)

func main() {
	log := log.New(os.Stdout, "Users: ", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)

	if err := run(log); err != nil {
		log.Printf("Main: %s", err)
		os.Exit(1)
	}
}

func run(log *log.Logger) error {

	ctx := context.WithValue(context.Background(), cmd.ContextKey(cmd.LOG), log)

	sess, err := uaws.GetSession(log)
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
