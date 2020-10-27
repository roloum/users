package main

import (
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

	"fmt"
	"os"
)

func main() {
	if err := run(); err != nil {
		fmt.Printf("Error: %s", err)
		os.Exit(1)
	}
}

func run() error {
	fmt.Println("Hello again")

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-west-2"),
	})
	if err != nil {
		return err
	}

	svc := dynamodb.New(sess)

	info := Info{"Rolando", "Umana"}
	user := User{"rolando.umana@gmail.com", time.Now().Format("2006-01-02"), info}

	av, err := dynamodbattribute.MarshalMap(user)
	if err != nil {
		return err
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String("User"),
	}

	_, err = svc.PutItem(input)
	if err != nil {
		return err
	}

	fmt.Println("Created...")

	return nil
}

//User ...
type User struct {
	Email   string `json:"email,omitempty"`
	Created string `json:"created,omitempty"`
	Info    Info   `json:"info,omitempty"`
}

//Info ...
type Info struct {
	FirstName string `json:"firstName,omitempty"`
	LastName  string `json:"lastName,omitempty"`
}
