package value_object

import "errors"

// We use VO as a Type Extension here
type Role string

const (
	Customer  Role = "customer"
	Organizer Role = "organizer"
	Admin     Role = "admin"
)

// NewRole validates and returns a new Role
func NewRole(r string) (Role, error) {
	if !isValidRole(r) {
		return Role(""), errors.New("invalid role")
	}
	return Role(r), nil
}

// isValidRole validates if new role value is valid
func isValidRole(role string) bool {
	switch role {
	case string(Customer), string(Organizer), string(Admin):
		return true
	default:
		return false
	}
}

// Value returns role of Role
func (r Role) Value() string {
	return string(r)
}

// EqualTo compares two Roles VO
func (r Role) EqualTo(other Role) bool {
	return r == other
}
