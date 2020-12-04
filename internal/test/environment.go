package test

import (
	"os"
)

//SetEnvironment sets the Environment variables for the test cacses
func SetEnvironment() {

	vars := map[string]string{
		"AWS_DYNAMODB_TABLE_USER": "User",
	}

	for key, value := range vars {
		_ = os.Setenv(key, value)
	}
}
