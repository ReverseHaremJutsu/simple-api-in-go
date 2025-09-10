package service

import (
	"rest-api-in-gin/internal/account/domain/model/value_object"
	"rest-api-in-gin/internal/account/domain/repository"
)

// UserAccountDomainService is a domain service
type UserAccountDomainService struct {
	userAccountChecker repository.UserAccountChecker
}

// NewUserAccountDomainService creates a new instance of UserAccountDomainService
func NewUserAccountDomainService(userAccountChecker repository.UserAccountChecker) *UserAccountDomainService {
	return &UserAccountDomainService{
		userAccountChecker: userAccountChecker,
	}
}

// IsEmailUnique validates if email VO already exist
func (s *UserAccountDomainService) IsEmailUnique(email value_object.Email) (bool, error) {
	exist, err := s.userAccountChecker.ExistByEmail(email)
	if err != nil {
		return false, err
	}
	return !exist, nil
}
