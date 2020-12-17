package user

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/aws/aws-lambda-go/events"
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

	//DynamoDBPrefixToken Prefix added to the sort key
	DynamoDBPrefixToken = "TOKEN"

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

//IsUserProfileKeys verifies that pk and sk correspond to a User's profile row
func IsUserProfileKeys(keys map[string]events.DynamoDBAttributeValue) bool {
	var userProfileKey bool

	log.Debug().Msgf("Checking User Profile keys")

	//Attribute pk exists in map?
	pk, ok := keys["pk"]
	if !ok {
		return false
	}

	//Attribute sk exists in map?
	sk, ok := keys["sk"]
	if !ok {
		return false
	}

	log.Debug().Msgf("Both keys are set")

	userProfileKey = strings.HasPrefix(pk.String(), DynamoDBPrefixUser) &&
		strings.HasPrefix(sk.String(), DynamoDBPrefixProfile)

	log.Debug().Msgf("IsUserProfileKeys: %v", userProfileKey)

	return userProfileKey
}

//Create creates a new user in DynamoDB and returns a pointer to the User object
//It inserts two rows in the table:
// - pk: USER#[email], sk: PROFILE# ... user profile row
// - pk: USER#[email], sk: TOKEN#[token] ... activation token (using id for now)
func Create(ctx context.Context, svc dynamodbiface.DynamoDBAPI, nu *NewUser,
	tableName string) (*User, error) {
	log.Info().Msgf("Creating user: %s", nu.Email)

	if tableName == "" {
		return nil, errors.New(ErrorUserTableNameIsEmpty)
	}

	log.Debug().Msg("Validating NewUser struct")

	if err := validate.Struct(nu); err != nil {
		return nil, getValidationError(err)
	}

	userID := uuid.New()
	log.Debug().Msgf("Generated UUID: %s", userID.String())

	u := &User{
		Email:     nu.Email,
		ID:        userID.String(),
		FirstName: nu.FirstName,
		LastName:  nu.LastName,
		Active:    false,
		Created:   time.Now().Format("2006-01-02"),
	}

	log.Debug().Msgf("Creating row: %+v", u)

	_, err := svc.TransactWriteItemsWithContext(ctx, &dynamodb.TransactWriteItemsInput{
		TransactItems: []*dynamodb.TransactWriteItem{
			{
				Put: &dynamodb.Put{
					Item: map[string]*dynamodb.AttributeValue{
						"pk":        {S: aws.String(u.getUserPK())},
						"sk":        {S: aws.String(u.getProfileSK())},
						"id":        {S: aws.String(u.ID)},
						"firstName": {S: aws.String(u.FirstName)},
						"lastName":  {S: aws.String(u.LastName)},
						"email":     {S: aws.String(u.Email)},
						"active":    {BOOL: aws.Bool(u.Active)},
						"created":   {S: aws.String(u.Created)},
						"type":      {S: aws.String(DynamoDBTypeUser)},
					},
					TableName:           aws.String(tableName),
					ConditionExpression: aws.String("attribute_not_exists(pk) and attribute_not_exists(sk)"),
				},
			},
			{
				Put: &dynamodb.Put{
					Item: map[string]*dynamodb.AttributeValue{
						"pk": {S: aws.String(u.getUserPK())},
						"sk": {S: aws.String(u.getTokenSK())},
					},
					TableName: aws.String(tableName),
				},
			},
		},
	})

	//	input := &dynamodb.PutItemInput{
	//		Item: map[string]*dynamodb.AttributeValue{
	//			"pk":        {S: aws.String(u.getUserPK())},
	//			"sk":        {S: aws.String(u.getProfileSK())},
	//			"id":        {S: aws.String(u.ID)},
	//			"firstName": {S: aws.String(u.FirstName)},
	//			"lastName":  {S: aws.String(u.LastName)},
	//			"email":     {S: aws.String(u.Email)},
	//			"active":    {BOOL: aws.Bool(u.Active)},
	//			"created":   {S: aws.String(u.Created)},
	//			"type":      {S: aws.String(DynamoDBTypeUser)},
	//		},
	//		ConditionExpression: aws.String("attribute_not_exists(pk) and attribute_not_exists(sk)"),
	//		TableName:           aws.String(tableName),
	//	}

	if err != nil {

		log.Debug().Msg(err.Error())

		if aerr, ok := err.(awserr.Error); ok {

			// TODO:
			//Error returned is TransactionCanceledException. Still trying to figure
			//out if there's a way to extract the Cancellation reasons
			if aerr.Code() == dynamodb.ErrCodeTransactionCanceledException {
				return nil, errors.New(ErrorDuplicateUser)
			}
		}

		return nil, err
	}

	return u, nil
}

//Activate sets the active column in the user-profile row to true
func (u *User) Activate(ctx context.Context, svc dynamodbiface.DynamoDBAPI,
	tableName string) error {

	log.Debug().Msgf("Activating user: %s", u.Email)

	return nil
}

func (u *User) getUserPK() string {
	return fmt.Sprintf("%s#%s", DynamoDBPrefixUser, u.Email)
}

func (u *User) getProfileSK() string {
	return fmt.Sprintf("%s#", DynamoDBPrefixProfile)
}

//getTokenSK forms Token SK with prefix and user ID
func (u *User) getTokenSK() string {
	return fmt.Sprintf("%s#%s", DynamoDBPrefixToken, u.ID)
}
