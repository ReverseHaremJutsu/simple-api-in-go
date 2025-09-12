package service

import (
	"rest-api-in-gin/internal/payment/application"
	"rest-api-in-gin/internal/payment/application/dto"
	"rest-api-in-gin/internal/payment/domain/repository"
	"rest-api-in-gin/internal/payment/domain/service"

	"github.com/google/uuid"
)

// DepositFundService is an application service to deposit fund into a wallet
type DepositFundService struct {
	walletWriter        repository.WalletWriter
	walletDomainService *service.WalletDomainService
}

// NewDepositFundService creates a new instance of DepositFundService
func NewDepositFundService(walletWriter repository.WalletWriter, walletDomainservice *service.WalletDomainService) *DepositFundService {
	return &DepositFundService{
		walletWriter:        walletWriter,
		walletDomainService: walletDomainservice,
	}
}

// Register consumes a request DTO, creates a payment and updates the wallet balance, returns a response DTO and error
func (s *DepositFundService) DepositFund(req *dto.DepositFundRequest) (*dto.DepositFundResponse, error) {

	exist, err := s.walletDomainService.IsExistingWallet(req.WalletID)
	if err != nil {
		return nil, application.NewAppError(application.ErrInternal, err.Error())
	}
	if !exist {
		return nil, application.NewAppError(application.ErrInvalidInput, "wallet does not exist")
	}

	updatedWallet, payment, err := s.walletDomainService.ReceivePayment(req.WalletID, uuid.Max, req.Amount)
	if err != nil {
		return nil, application.NewAppError(application.ErrInternal, "domain error receiving payment")
	}

	updates := []repository.WalletPaymentUpdate{
		{Wallet: updatedWallet, Payment: payment},
	}

	if err := s.walletWriter.UpdatesByID(updates); err != nil {
		return nil, application.NewAppError(application.ErrInternal, "repository error receiving payment")
	}

	return dto.NewDepositFundResponse(req.WalletID, req.Amount, updatedWallet.Balance.Value()), nil
}
