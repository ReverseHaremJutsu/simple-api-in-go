package dto

import (
	"rest-api-in-gin/internal/account/domain/model/entity"

	"github.com/google/uuid"
)

// UserAccountDTO is from application to client, for returning responses
type UserAccountDTO struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Email string    `json:"email"`
}

// NewUserAccountDTO maps a domain UserAccount entity to UserDTO
func NewUserAccountDTO(user *entity.UserAccount) *UserAccountDTO {
	return &UserAccountDTO{
		ID:    user.ID,
		Name:  user.Name.Value(),
		Email: user.Email.Value(),
	}
}
