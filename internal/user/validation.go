package user

import (
	"errors"

	validator "github.com/go-playground/validator/v10"
	emailaddress "github.com/mcnijman/go-emailaddress"
)

const (
	//ErrorFirstNameIsEmpty Error describes first name being empty
	ErrorFirstNameIsEmpty = "FirstNameIsEmpty"

	//ErrorLastNameIsEmpty Error describes last name being empty
	ErrorLastNameIsEmpty = "LastNameIsEmpty"

	//ErrorEmailIsEmpty Error describes email being empty
	ErrorEmailIsEmpty = "EmailIsEmpty"

	//ErrorInvalidEmail Error describes email being invalid
	ErrorInvalidEmail = "InvalidEmail"
)

var validate *validator.Validate

//init instantiates a validator
func init() {
	validate = validator.New()
}

//getValidationError Returns the first error reported by the validator
func getValidationError(verr error) error {

	//Retrieve first error
	err := verr.(validator.ValidationErrors)[0]

	switch err.Field() {
	case "FirstName":
		return errors.New(ErrorFirstNameIsEmpty)
	case "LastName":
		return errors.New(ErrorLastNameIsEmpty)
	case "Email":
		switch err.Tag() {
		case "required":
			return errors.New(ErrorEmailIsEmpty)
		case "email":
			return errors.New(ErrorInvalidEmail)
		}
	}

	return nil
}

//isValidEmail validates an email address
func isValidEmail(email string) error {
	if _, err := emailaddress.Parse(email); err != nil {

		return errors.New(ErrorInvalidEmail)
	}
	return nil
}
