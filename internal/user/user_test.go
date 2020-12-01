package user

import (
	"context"
	"errors"
	"io/ioutil"
	"log"
	"reflect"
	"testing"

	"github.com/roloum/users/internal/test"
)

//TestCreateUser Tests the Create functionality
func TestCreateUser(t *testing.T) {

	tests := []struct {
		desc      string
		user      *NewUser
		mockDBSvc *test.MockDynamoDB
		err       error
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
		},
		{
			desc: ErrorFirstNameIsEmpty,
			user: &NewUser{
				LastName: "User",
				Email:    "test@user.com",
			},
			mockDBSvc: &test.MockDynamoDB{},
			err:       errors.New(ErrorFirstNameIsEmpty),
		},
		{
			desc: ErrorLastNameIsEmpty,
			user: &NewUser{
				FirstName: "Test",
				Email:     "test@user.com",
			},
			mockDBSvc: &test.MockDynamoDB{},
			err:       errors.New(ErrorLastNameIsEmpty),
		},
		{
			desc: ErrorEmailIsEmpty,
			user: &NewUser{
				FirstName: "Test",
				LastName:  "User",
			},
			mockDBSvc: &test.MockDynamoDB{},
			err:       errors.New(ErrorEmailIsEmpty),
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
		},
	}

	//Do not log
	log := log.New(ioutil.Discard, "", 0)

	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			_, err := Create(context.Background(), tc.mockDBSvc, tc.user, log)
			if !reflect.DeepEqual(err, tc.err) {
				t.Errorf("Expected: %v. Received: %v", err, tc.err)
			}
		})
	}

}
