// filepath: /home/tshering/Desktop/realworld-example-app/golang-gin-realworld-example-app/users/validators.go

package users

import (
	"errors"
	"regexp"
)

// UserValidator validates user input for registration and login
type UserValidator struct{}

// ValidateRegistration validates the user registration input
func (uv *UserValidator) ValidateRegistration(username, email, password string) error {
	if len(username) < 3 {
		return errors.New("username must be at least 3 characters long")
	}
	if !isValidEmail(email) {
		return errors.New("invalid email format")
	}
	if len(password) < 6 {
		return errors.New("password must be at least 6 characters long")
	}
	return nil
}

// ValidateLogin validates the user login input
func (uv *UserValidator) ValidateLogin(email, password string) error {
	if !isValidEmail(email) {
		return errors.New("invalid email format")
	}
	if len(password) == 0 {
		return errors.New("password cannot be empty")
	}
	return nil
}

// isValidEmail checks if the email format is valid
func isValidEmail(email string) bool {
	// Simple regex for validating an email
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}