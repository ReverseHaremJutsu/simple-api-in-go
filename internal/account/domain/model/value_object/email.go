package value_object

import (
	"errors"
	"regexp"
	"strings"
)

// Email is a Value Object
type Email struct {
	emailAddress string
}

// NewEmail validates and returns a new Email VO
func NewEmail(email string) (Email, error) {
	email = strings.TrimSpace(email)
	if !isValidEmail(email) {
		return Email{}, errors.New("invalid email format")
	}
	return Email{emailAddress: email}, nil
}

// isValidEmail checks if new emailAddress value is valid
func isValidEmail(email string) bool {
	// Simple regex to validate an Email (for now)
	var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

// Value returns emailAddress of Email
func (e Email) Value() string {
	return e.emailAddress
}

// EqualTo compares two EmailAddress VO
func (e Email) EqualTo(other Email) bool {
	return e.emailAddress == other.emailAddress
}
