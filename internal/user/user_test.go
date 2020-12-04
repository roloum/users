package user

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"reflect"
	"testing"

	"github.com/roloum/users/internal/config"
	"github.com/roloum/users/internal/test"
)

func init() {
	fmt.Println("init user user_test.go")
	test.SetEnvironment()
}

//TestCreateUser Tests the Create functionality
func TestCreateUser(t *testing.T) {

	var cfg config.Config
	_ = cfg

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

	//Do not log
	log := log.New(ioutil.Discard, "", 0)

	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			//"User"
			_, err := Create(context.Background(), tc.mockDBSvc, tc.user, tc.tableName,
				log)
			if !reflect.DeepEqual(err, tc.err) {
				t.Errorf("Expected: %v. Received: %v", tc.err, err)
			}
		})
	}

}
