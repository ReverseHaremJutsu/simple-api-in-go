package repository

import (
	"rest-api-in-gin/internal/account/domain/model/entity"
	"rest-api-in-gin/internal/account/domain/model/value_object"

	"github.com/google/uuid"
)

// UserAccountReader is a interface contract
type UserAccountReader interface {
	// TODO
	GetByID(id uuid.UUID) (*entity.UserAccount, error)
	GetByEmail(email value_object.Email) (*entity.UserAccount, error)
	GetAllByRole(role value_object.Role) ([]*entity.UserAccount, error)
}
