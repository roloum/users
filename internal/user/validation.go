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

	validate.RegisterValidation("validEmail", isValidEmail)
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
			//		case "email":
			//			return errors.New(ErrorInvalidEmail)
		case "validEmail":
			return errors.New(ErrorInvalidEmail)
		}
	}
	return nil
}

//isValidEmail validates an email address using go-emailaddress
//From validator: This validates that a string value contains a valid email
//This may not conform to all possibilities of any rfc standard, but neither
//does any email provider accept all possibilities.
func isValidEmail(fl validator.FieldLevel) bool {
	if _, err := emailaddress.Parse(fl.Field().String()); err != nil {
		return false
	}
	return true
}
