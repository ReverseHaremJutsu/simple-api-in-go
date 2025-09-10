package repository

import (
	"rest-api-in-gin/internal/account/domain/model/value_object"
)

// UserAccountChecker is a interface contract
type UserAccountChecker interface {
	ExistByEmail(email value_object.Email) (bool, error)
	// TODO
	// ExistByID(id uuid.UUID) (bool, error)
}
