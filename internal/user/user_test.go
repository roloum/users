package user

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/aws/aws-lambda-go/events"

	"github.com/roloum/users/internal/test"
	"github.com/rs/zerolog"
)

func init() {
	zerolog.SetGlobalLevel(zerolog.Disabled)
	test.SetEnvironment()
}

//TestCreateUser Tests the Create functionality
func TestCreateUser(t *testing.T) {

	tests := []struct {
		desc      string
		user      *NewUser
		mockDBSvc *test.MockDynamoDB
		err       error
		tableName string
	}{
		{
			desc: "CreateUser",
			user: &NewUser{
				FirstName: "Test",
				LastName:  "User",
				Email:     "test@user.com",
			},
			mockDBSvc: &test.MockDynamoDB{},
			err:       nil,
			tableName: "User",
		},
		{
			desc: ErrorFirstNameIsEmpty,
			user: &NewUser{
				LastName: "User",
				Email:    "test@user.com",
			},
			mockDBSvc: &test.MockDynamoDB{},
			err:       errors.New(ErrorFirstNameIsEmpty),
			tableName: "User",
		},
		{
			desc: ErrorLastNameIsEmpty,
			user: &NewUser{
				FirstName: "Test",
				Email:     "test@user.com",
			},
			mockDBSvc: &test.MockDynamoDB{},
			err:       errors.New(ErrorLastNameIsEmpty),
			tableName: "User",
		},
		{
			desc: ErrorEmailIsEmpty,
			user: &NewUser{
				FirstName: "Test",
				LastName:  "User",
			},
			mockDBSvc: &test.MockDynamoDB{},
			err:       errors.New(ErrorEmailIsEmpty),
			tableName: "User",
		},
		{
			desc: ErrorInvalidEmail,
			user: &NewUser{
				FirstName: "Test",
				LastName:  "User",
				Email:     "yadayadayada",
			},
			mockDBSvc: &test.MockDynamoDB{},
			err:       errors.New(ErrorInvalidEmail),
			tableName: "User",
		},
		{
			desc: ErrorUserTableNameIsEmpty,
			user: &NewUser{
				FirstName: "Test",
				LastName:  "User",
				Email:     "yadayadayada",
			},
			mockDBSvc: &test.MockDynamoDB{},
			err:       errors.New(ErrorUserTableNameIsEmpty),
			tableName: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			//"User"
			_, err := Create(context.Background(), tc.mockDBSvc, tc.user, tc.tableName)
			if !reflect.DeepEqual(err, tc.err) {
				t.Errorf("Expected: %v. Received: %v", tc.err, err)
			}
		})
	}

}

func TestIsUserProfileKeys(t *testing.T) {
	tests := []struct {
		desc     string
		keys     map[string]events.DynamoDBAttributeValue
		expected bool
	}{
		{
			desc: "profileKeys",
			keys: map[string]events.DynamoDBAttributeValue{
				"pk": events.NewStringAttribute("USER#"),
				"sk": events.NewStringAttribute("PROFILE#"),
			},
			expected: true,
		},
		{
			desc: "differentKeys",
			keys: map[string]events.DynamoDBAttributeValue{
				"pk": events.NewStringAttribute("USER#"),
				"sk": events.NewStringAttribute("DIFF#"),
			},
			expected: false,
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			if result := IsUserProfileKeys(tc.keys); result != tc.expected {
				t.Errorf("Expected: %v. Received: %v", tc.expected, result)
			}
		})
	}
}
