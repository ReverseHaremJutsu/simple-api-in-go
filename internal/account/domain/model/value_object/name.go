package value_object

import (
	"errors"
	"regexp"
	"strings"
)

// Name is a Value Object
type Name struct {
	name string
}

// NewName validates and returns a new Name VO
func NewName(name string) (Name, error) {
	name = strings.TrimSpace(name)
	if !isValidName(name) {
		return Name{}, errors.New("invalid name format")
	}
	return Name{name: name}, nil
}

// isValidName checks if new name value is valid
func isValidName(name string) bool {

	var nameRegex = regexp.MustCompile(`^[A-Za-z ]+$`)
	return nameRegex.MatchString(name)
}

// Value returns name of Name
func (n Name) Value() string {
	return n.name
}

// EqualTo compares two Name VO
func (n Name) EqualTo(other Name) bool {
	return n.name == other.name
}
