package service

import (
	"rest-api-in-gin/internal/account/application"
	"rest-api-in-gin/internal/account/application/dto"
	"rest-api-in-gin/internal/account/domain/model/entity"
	"rest-api-in-gin/internal/account/domain/model/value_object"
	"rest-api-in-gin/internal/account/domain/repository"
	"rest-api-in-gin/internal/account/domain/service"
)

// RegisterAccountService is an application service to register a new UserAccount
type RegisterAccountService struct {
	userAccountWriter        repository.UserAccountWriter
	userAccountDomainService *service.UserAccountDomainService
}

// NewRegisterAccountService creates a new instance of RegisterAccountService
func NewRegisterAccountService(userAccountWriter repository.UserAccountWriter, userAccountDomainService *service.UserAccountDomainService) *RegisterAccountService {
	return &RegisterAccountService{
		userAccountWriter:        userAccountWriter,
		userAccountDomainService: userAccountDomainService,
	}
}

// Register consumes a RegisterDTO, creates and persist a new UserAccountand and, returns a UserAccountDTO and error
func (s *RegisterAccountService) Register(req *dto.RegisterDTO) (*dto.UserAccountDTO, error) {

	newemail, err := value_object.NewEmail(req.Email)
	if err != nil {
		return nil, application.NewAppError(application.ErrInvalidInput, err.Error())
	}

	unique, err := s.userAccountDomainService.IsEmailUnique(newemail)
	if err != nil {
		return nil, application.NewAppError(application.ErrInternal, err.Error())
	}
	if !unique {
		return nil, application.NewAppError(application.ErrInvalidInput, "email already exist")
	}

	newuseraccount, err := entity.NewUserAccount(
		req.Name,
		req.Email,
		req.Password,
		req.Role,
	)
	if err != nil {
		return nil, application.NewAppError(application.ErrInvalidInput, err.Error())
	}

	if err := s.userAccountWriter.Create(newuseraccount); err != nil {
		return nil, application.NewAppError(application.ErrInternal, err.Error())
	}

	return dto.NewUserAccountDTO(newuseraccount), nil
}
