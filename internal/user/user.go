package user

import (
	"context"
	"errors"
	"log"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
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
	Email     string `json:"email" validate:"required,validEmail"`
	FirstName string `json:"firstName" validate:"required"`
	LastName  string `json:"lastName" validate:"required"`
}

//Create creates a new user in DynamoDB and returns a pointer to the User
//object
func Create(ctx context.Context, dynamoDB dynamodbiface.DynamoDBAPI, nu *NewUser,
	tableName string, log *log.Logger) (*User, error) {
	log.Printf("Creating user: %v\n", nu)

	if tableName == "" {
		return nil, errors.New(ErrorUserTableNameIsEmpty)
	}

	if err := validate.Struct(nu); err != nil {
		return nil, getValidationError(err)
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
		TableName:           aws.String(tableName),
	}

	if _, err := dynamoDB.PutItemWithContext(ctx, input); err != nil {
		return nil, err
	}

	log.Printf("User created: %v\n", *u)
	return u, nil
}
