package repository

import (
	"rest-api-in-gin/internal/account/domain/model/entity"
)

// UserAccountWriter is a interface contract
type UserAccountWriter interface {
	Create(useraccount *entity.UserAccount) error
	// TODO
	// UpdateByID(id uuid.UUID, useraccount *entity.UserAccount) error
	// DeleteByID(id uuid.UUID) error
}
