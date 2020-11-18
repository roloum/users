package user

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	emailaddress "github.com/mcnijman/go-emailaddress"
)

//User contains information about the user
type User struct {
	ID        string `json:"id,omitempty"`
	Email     string `json:"email,omitempty"`
	FirstName string `json:"firstName,omitempty"`
	LastName  string `json:"lastName,omitempty"`
	Active    bool   `json:"active,omitempty"`
	Created   string `json:"created,omitempty"`
}

//NewUser contains information to create new user
type NewUser struct {
	Email     string `json:"email" validate:"required,email"`
	FirstName string `json:"firstName" validate:"required"`
	LastName  string `json:"lastName" validate:"required"`
}

//Create creates a new user in DynamoDB and returns a pointer to the User
//object
func Create(ctx context.Context, dynamoDB *dynamodb.DynamoDB, nu *NewUser,
	log *log.Logger) (*User, error) {
	log.Printf("Creating user: %v\n", nu)

	if err := validateNew(nu); err != nil {
		log.Printf("Error validating new user: %s", err)
		return nil, err
	}

	u := &User{
		Email:     strings.ToLower(nu.Email),
		FirstName: nu.FirstName,
		LastName:  nu.LastName,
		Active:    false,
		Created:   time.Now().Format("2006-01-02"),
	}

	input := &dynamodb.PutItemInput{
		Item: map[string]*dynamodb.AttributeValue{
			"email":     {S: aws.String(u.Email)},
			"created":   {S: aws.String(u.Created)},
			"firstName": {S: aws.String(u.FirstName)},
			"lastName":  {S: aws.String(u.LastName)},
			"active":    {BOOL: aws.Bool(u.Active)},
		},
		ConditionExpression: aws.String("attribute_not_exists(email)"),
		TableName:           aws.String("User"),
	}

	if _, err := dynamoDB.PutItemWithContext(ctx, input); err != nil {
		return nil, err
	}

	log.Printf("User created: %v\n", *u)
	return u, nil
}

//validateNew validates the new user
func validateNew(nu *NewUser) error {
	if err := isValidEmail(nu.Email); err != nil {
		return err
	}
	if nu.FirstName == "" {
		return errors.New("First name is empty")
	}
	if nu.LastName == "" {
		return errors.New("First name is empty")
	}
	return nil
}

//isValidEmail validates an email address
func isValidEmail(email string) error {
	if _, err := emailaddress.Parse(email); err != nil {

		return fmt.Errorf("Email: %s", err)
	}
	return nil
}
