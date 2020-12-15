package user

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"

	"github.com/google/uuid"
)

const (
	//DynamoDBPrefixUser Prexix added to the primary key
	DynamoDBPrefixUser = "USER"

	//DynamoDBPrefixProfile Prefix added to the sort key
	DynamoDBPrefixProfile = "PROFILE"

	//DynamoDBTypeUser identifies the type of row in dynamoDB
	DynamoDBTypeUser = "User"

	//ErrorDuplicateUser Returned when the user already exists in the table
	ErrorDuplicateUser = "DuplicatedUser"

	//ErrorUserTableNameIsEmpty Error describes AWS table name being empty
	ErrorUserTableNameIsEmpty = "UserTableNameIsEmpty"
)

//User contains information about the user
type User struct {
	ID        string `json:"id,omitempty"`
	FirstName string `json:"firstName,omitempty"`
	LastName  string `json:"lastName,omitempty"`
	Email     string `json:"email,omitempty"`
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

	log.Println("Validation NewUser struct")

	if err := validate.Struct(nu); err != nil {
		return nil, getValidationError(err)
	}

	userID := uuid.New()
	log.Printf("Generated UUID: %s", userID.String())

	u := &User{
		Email:     nu.Email,
		ID:        userID.String(),
		FirstName: nu.FirstName,
		LastName:  nu.LastName,
		Active:    false,
		Created:   time.Now().Format("2006-01-02"),
	}

	log.Printf("Creating row: %+v", u)

	input := &dynamodb.PutItemInput{
		Item: map[string]*dynamodb.AttributeValue{
			"pk":        {S: aws.String(u.getPK())},
			"sk":        {S: aws.String(u.getSK())},
			"id":        {S: aws.String(u.ID)},
			"firstName": {S: aws.String(u.FirstName)},
			"lastName":  {S: aws.String(u.LastName)},
			"email":     {S: aws.String(u.Email)},
			"active":    {BOOL: aws.Bool(u.Active)},
			"created":   {S: aws.String(u.Created)},
			"type":      {S: aws.String(DynamoDBTypeUser)},
		},
		ConditionExpression: aws.String("attribute_not_exists(pk) and attribute_not_exists(sk)"),
		TableName:           aws.String(tableName),
	}

	if _, err := dynamoDB.PutItemWithContext(ctx, input); err != nil {

		fmt.Printf("Error: %v", err)
		if aerr, ok := err.(awserr.Error); ok {
			if aerr.Code() == dynamodb.ErrCodeConditionalCheckFailedException {
				return nil, errors.New(ErrorDuplicateUser)
			}
		}

		return nil, err
	}

	log.Printf("User created: %v\n", *u)
	return u, nil
}

func (u *User) getPK() string {
	return fmt.Sprintf("%s#%s", DynamoDBPrefixUser, u.Email)
}

func (u *User) getSK() string {
	return fmt.Sprintf("%s#", DynamoDBPrefixProfile)
}
