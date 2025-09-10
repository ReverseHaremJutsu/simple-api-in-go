package entity

import (
	"fmt"
	"time"

	"rest-api-in-gin/internal/account/domain/model/value_object"

	"github.com/google/uuid"
)

// UserAccount is an Aggregate Entity
type UserAccount struct {
	ID        uuid.UUID
	Name      value_object.Name
	Email     value_object.Email
	Password  value_object.Password
	Role      value_object.Role
	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewUserAccount creates a new UserAccount
func NewUserAccount(name string, email string, rawPassword string, role string) (*UserAccount, error) {

	nameVO, err := value_object.NewName(name)
	if err != nil {
		return nil, fmt.Errorf("error creating account: %w", err)
	}

	emailVO, err := value_object.NewEmail(email)
	if err != nil {
		return nil, fmt.Errorf("error creating account: %w", err)
	}

	passwordVO, err := value_object.NewPassword(rawPassword)
	if err != nil {
		return nil, fmt.Errorf("error creating account: %w", err)
	}

	roleVO, err := value_object.NewRole(role)
	if err != nil {
		return nil, fmt.Errorf("error creating account: %w", err)
	}

	now := time.Now()

	return &UserAccount{
		ID:        uuid.New(),
		Name:      nameVO,
		Email:     emailVO,
		Password:  passwordVO,
		Role:      roleVO,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

// ChangeName updates the UserAccount name
func (u *UserAccount) ChangeName(newName string) error {
	name, err := value_object.NewName(newName)
	if err != nil {
		return err
	}
	u.Name = name
	u.UpdatedAt = time.Now()
	return nil
}

// ChangePassword updates the UserAccount password
func (u *UserAccount) ChangePassword(newPassword string) error {
	password, err := value_object.NewPassword(newPassword)
	if err != nil {
		return err
	}
	u.Password = password
	u.UpdatedAt = time.Now()
	return nil
}
