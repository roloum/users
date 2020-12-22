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

const (
	UserTable = "User"
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
			tableName: UserTable,
		},
		{
			desc: ErrorFirstNameIsEmpty,
			user: &NewUser{
				LastName: "User",
				Email:    "test@user.com",
			},
			mockDBSvc: &test.MockDynamoDB{},
			err:       errors.New(ErrorFirstNameIsEmpty),
			tableName: UserTable,
		},
		{
			desc: ErrorLastNameIsEmpty,
			user: &NewUser{
				FirstName: "Test",
				Email:     "test@user.com",
			},
			mockDBSvc: &test.MockDynamoDB{},
			err:       errors.New(ErrorLastNameIsEmpty),
			tableName: UserTable,
		},
		{
			desc: ErrorEmailIsEmpty,
			user: &NewUser{
				FirstName: "Test",
				LastName:  "User",
			},
			mockDBSvc: &test.MockDynamoDB{},
			err:       errors.New(ErrorEmailIsEmpty),
			tableName: UserTable,
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
			tableName: UserTable,
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

	t.Run("profileKeys", func(t *testing.T) {
		keys := map[string]events.DynamoDBAttributeValue{
			"pk": events.NewStringAttribute("USER#"),
			"sk": events.NewStringAttribute("PROFILE#"),
		}
		if result := IsUserProfileKeys(keys); !result {
			t.Errorf("Expected: %v.", result)
		}
	})

	t.Run("tokenKeys", func(t *testing.T) {
		keys := map[string]events.DynamoDBAttributeValue{
			"pk": events.NewStringAttribute("USER#"),
			"sk": events.NewStringAttribute("TOKEN#"),
		}
		if result := IsUserTokenKeys(keys); !result {
			t.Errorf("Expected: %v.", result)
		}
	})

	t.Run("differentKeys", func(t *testing.T) {
		keys := map[string]events.DynamoDBAttributeValue{
			"pk": events.NewStringAttribute("USER#"),
			"sk": events.NewStringAttribute("DIFF#"),
		}
		if result := IsUserProfileKeys(keys); result {
			t.Errorf("Expected: %v.", !result)
		}
	})
}
