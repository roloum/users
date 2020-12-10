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
)

func initHandler(ctx context.Context, e events.DynamoDBEvent) {

	log := log.New(os.Stdout, "notifyUser: ", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)
	log.Printf("%+v", e)

}

func main() {
	lambda.Start(initHandler)
}
