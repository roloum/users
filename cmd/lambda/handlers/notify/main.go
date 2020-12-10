//Lambda function that sends an email to the user after:
// - user is created
// - user is verified
package main

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	uaws "github.com/roloum/users/internal/aws"
	"github.com/roloum/users/internal/user"
)

func handler(ctx context.Context, e events.DynamoDBEvent, log *log.Logger) {

	for _, v := range e.Records {
		log.Printf("Event name: %s\n", v.EventName)

		switch v.EventName {
		case "INSERT":

			var u user.User
			err := uaws.UnmarshalStreamImage(v.Change.NewImage, &u)
			if err != nil {
				log.Fatal(err)
			}
			log.Println(v.Change.NewImage)
			sendEmail(
				"Activate account",
				"body",
				u.Email,
				log)
			break
		case "MODIFY":

			log.Println(v.Change.OldImage)
			sendEmail(
				"Welcome!",
				"Welcome",
				"email",
				log)
		}
	}
}

func sendEmail(subject, body, email string, log *log.Logger) error {
	log.Printf("Sending email:\nemail: %s\n, Subject: %s\nBody: %s\n",
		email, subject, body)
	return nil
}

func initHandler(ctx context.Context, e events.DynamoDBEvent) {

	log := log.New(os.Stdout, "notifyUser: ", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)
	log.Printf("%+v", e)

	handler(ctx, e, log)

}

func main() {
	lambda.Start(initHandler)
}
