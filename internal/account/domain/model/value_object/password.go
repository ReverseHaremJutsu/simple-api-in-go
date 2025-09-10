package value_object

import (
	"errors"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

// Password is a Value Object
type Password struct {
	hashedPassword string
}

// NewPassword validates, hashes the raw password and, returns a new Password VO
func NewPassword(rawPassword string) (Password, error) {
	rawPassword = strings.TrimSpace(rawPassword)
	if !isValidRawPassword(rawPassword) {
		return Password{}, errors.New("invalid password format")
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte(rawPassword), bcrypt.DefaultCost)
	if err != nil {
		return Password{}, errors.New("failed to hash password")
	}
	return Password{hashedPassword: string(hashed)}, nil
}

// isValidRawPassword validates if new raw password value is valid
func isValidRawPassword(rawPassword string) bool {
	return len(rawPassword) >= 4
}

// Value returns the hashed password string
func (p Password) Value() string {
	return p.hashedPassword
}

// Matches compares a Password VO with another raw password string
// Deviate slightly from rule of validation of VO equality
func (p Password) Matches(otherRawPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(p.hashedPassword), []byte(otherRawPassword))
	return err == nil
}
