//Lambda function that sends an email to the user after:
// - user is created
// - user is verified
package main

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/rs/zerolog/log"

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
	sender string) error {

	for _, v := range e.Records {
		log.Debug().Msgf("Event name: %s\n", v.EventName)

		hostname, err := os.Hostname()
		if err != nil {
			return err
		}

		switch v.EventName {
		case "INSERT":

			var u user.User

			log.Debug().Msg("Unmarshalling user struct")

			//Unmarshal Image into user struct
			err := uaws.UnmarshalStreamImage(v.Change.NewImage, &u)
			if err != nil {
				log.Fatal().Msg(err.Error())
			}

			activationURL := fmt.Sprintf("%s/dev/users/activate", hostname)

			if err := sendEmail(
				"Activate account",
				fmt.Sprintf("<a href=\"%s\">Click here to Activate</a>", activationURL),
				fmt.Sprintf("Click here to Activate: \"%s\"", activationURL),
				u.Email,
				sender,
				svc); err != nil {
				log.Fatal().Msg(err.Error())
			}
		case "MODIFY":
		}
	}

	return nil
}

func sendEmail(subject, htmlBody, textBody, recipient, sender string,
	svc *ses.SES) error {
	log.Debug().Str("sender", sender).
		Str("nemail", recipient).
		Str("subject", subject).
		Str("body", textBody).
		Msg("Sending email")

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

	log.Debug().Msgf("%+v", result)
	log.Info().Msg("Email sent")

	return nil
}

func initHandler(ctx context.Context, e events.DynamoDBEvent) error {

	//Config holds the configuration for the application
	var cfg configuration
	err := config.Load(&cfg)
	if err != nil {
		return err
	}

	sess, err := uaws.GetSession(cfg.AWS.Region)
	if err != nil {
		return err
	}

	return handler(ctx, e, ses.New(sess), cfg.Email.Sender)

}

func main() {
	lambda.Start(initHandler)
}
