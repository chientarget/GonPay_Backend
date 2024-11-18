// pkg/validator/validator.go
package validator

import (
	"fmt"
	"regexp"
	"strings"
)

type ValidatorInterface interface {
	ValidateEmail(email string) error
	ValidatePassword(password string) error
	ValidatePhone(phone string) error
	ValidateUsername(username string) error
	ValidateAmount(amount float64) error
}

type Validator struct{}

func NewValidator() ValidatorInterface {
	return &Validator{}
}

func (v *Validator) ValidateEmail(email string) error {
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	if !emailRegex.MatchString(strings.ToLower(email)) {
		return fmt.Errorf("invalid email format")
	}
	return nil
}

func (v *Validator) ValidatePassword(password string) error {
	if len(password) < 6 {
		return fmt.Errorf("password must be at least 6 characters long")
	}

	hasUpperCase := false
	hasLowerCase := false
	hasNumber := false

	for _, char := range password {
		switch {
		case 'a' <= char && char <= 'z':
			hasLowerCase = true
		case 'A' <= char && char <= 'Z':
			hasUpperCase = true
		case '0' <= char && char <= '9':
			hasNumber = true
		}
	}

	if !hasUpperCase || !hasLowerCase || !hasNumber {
		return fmt.Errorf("password must contain at least one uppercase letter, one lowercase letter, and one number")
	}

	return nil
}

func (v *Validator) ValidatePhone(phone string) error {
	phoneRegex := regexp.MustCompile(`^\+?[1-9]\d{1,14}$`)
	if !phoneRegex.MatchString(phone) {
		return fmt.Errorf("invalid phone number format")
	}
	return nil
}

func (v *Validator) ValidateUsername(username string) error {
	if len(username) < 3 || len(username) > 50 {
		return fmt.Errorf("username must be between 3 and 50 characters")
	}
	usernameRegex := regexp.MustCompile(`^[a-zA-Z0-9_\-]+$`)
	if !usernameRegex.MatchString(username) {
		return fmt.Errorf("username can only contain letters, numbers, underscores and hyphens")
	}
	return nil
}

func (v *Validator) ValidateAmount(amount float64) error {
	if amount <= 0 {
		return fmt.Errorf("amount must be greater than 0")
	}
	return nil
}
