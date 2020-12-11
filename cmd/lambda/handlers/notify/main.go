//Lambda function that sends an email to the user after:
// - user is created
// - user is verified
package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ses"
	uaws "github.com/roloum/users/internal/aws"
	"github.com/roloum/users/internal/config"
	"github.com/roloum/users/internal/user"
)

//CHARSET Character encoding for email
const CHARSET = "UTF-8"

type configuration struct {
	AWS struct {
		Region string `required:"true"`
	}
	Email struct {
		Sender string `required:"true"`
	}
}

func handler(ctx context.Context, e events.DynamoDBEvent, svc *ses.SES,
	sender string, log *log.Logger) error {

	for _, v := range e.Records {
		log.Printf("Event name: %s\n", v.EventName)

		hostname, err := os.Hostname()
		if err != nil {
			return err
		}

		switch v.EventName {
		case "INSERT":

			var u user.User
			//Unmarshal Image into user struct
			err := uaws.UnmarshalStreamImage(v.Change.NewImage, &u)
			if err != nil {
				log.Fatal(err)
			}

			activationURL := fmt.Sprintf("%s/dev/users/activate", hostname)

			if err := sendEmail(
				"Activate account",
				fmt.Sprintf("<a href=\"%s\">Click here to Activate</a>", activationURL),
				fmt.Sprintf("Click here to Activate: \"%s\"", activationURL),
				u.Email,
				sender,
				svc,
				log); err != nil {
				log.Fatal(err)
			}
		case "MODIFY":
		}
	}

	return nil
}

func sendEmail(subject, htmlBody, textBody, recipient, sender string,
	svc *ses.SES, log *log.Logger) error {
	log.Printf("Sending email:\nsender: %s, email: %s,\nSubject: %s,\nBody: %s\n",
		sender, recipient, subject, textBody)

	// Email input
	input := &ses.SendEmailInput{
		Destination: &ses.Destination{
			ToAddresses: []*string{
				aws.String(recipient),
			},
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Html: &ses.Content{
					Charset: aws.String(CHARSET),
					Data:    aws.String(htmlBody),
				},
				Text: &ses.Content{
					Charset: aws.String(CHARSET),
					Data:    aws.String(textBody),
				},
			},
			Subject: &ses.Content{
				Charset: aws.String(CHARSET),
				Data:    aws.String(subject),
			},
		},
		Source: aws.String(sender),
	}

	// Attempt to send the email.
	result, err := svc.SendEmail(input)

	// Display error messages if they occur.
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			return errors.New(aerr.Code())
		}
		return err
	}

	log.Println("Email sent")
	log.Println(result)

	return nil
}

func initHandler(ctx context.Context, e events.DynamoDBEvent) error {

	log := log.New(os.Stdout, "notifyUser: ", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)

	//Config holds the configuration for the application
	var cfg configuration
	err := config.Load(&cfg, log)
	if err != nil {
		return err
	}

	sess, err := uaws.GetSession(cfg.AWS.Region, log)
	if err != nil {
		return err
	}

	return handler(ctx, e, ses.New(sess), cfg.Email.Sender, log)

}

func main() {
	lambda.Start(initHandler)
}
