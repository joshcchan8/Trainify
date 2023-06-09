package models

import (
	"errors"
	"regexp"
)

type User struct {
	UserID    int    `json:"user_id"`    // required
	Username  string `json:"username"`   // required
	Email     string `json:"email"`      // required, verified with regex
	Password  string `json:"password"`   // required, verified with regex (unguessable), will be hashed
	ProfileID int    `json:"profile_id"` // required
}

type Info struct {
	Username string
	Email    string
	Password string
}

func ValidateEmail(email string) error {

	var invalidErr error = nil

	// Regular expression pattern for email validation
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	match, _ := regexp.MatchString(pattern, email)

	if !match {
		invalidErr = errors.New("invalid email address format")
	}

	return invalidErr
}

// Returns an error if there is insufficient registration information
func RegistrationInfoProvided(registrationInfo Info) error {

	var missingInfoErr error = nil

	if registrationInfo.Username == "" {
		missingInfoErr = errors.New("missing username")
	} else if registrationInfo.Email == "" {
		missingInfoErr = errors.New("missing email")
	} else if registrationInfo.Password == "" {
		missingInfoErr = errors.New("missing password")
	}

	return missingInfoErr
}

// Returns an error if there is insufficient login information
func LoginInfoProvided(loginInfo Info) error {

	var missingInfoErr error = nil

	if loginInfo.Email == "" {
		missingInfoErr = errors.New("missing email")
	} else if loginInfo.Password == "" {
		missingInfoErr = errors.New("missing password")
	}

	return missingInfoErr
}
