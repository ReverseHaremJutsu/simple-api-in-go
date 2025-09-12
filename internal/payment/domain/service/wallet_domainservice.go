package service

import (
	"rest-api-in-gin/internal/payment/domain/model/entity"

	"rest-api-in-gin/internal/payment/domain/repository"

	"github.com/google/uuid"
)

// WalletDomainService is a domain service
type WalletDomainService struct {
	walletChecker repository.WalletChecker
	walletReader  repository.WalletReader
	walletWriter  repository.WalletWriter
}

// NewWalletDomainService creates a new instance of WalletDomainService
func NewWalletDomainService(walletChecker repository.WalletChecker, walletReader repository.WalletReader, walletWriter repository.WalletWriter) *WalletDomainService {
	return &WalletDomainService{
		walletChecker: walletChecker,
		walletReader:  walletReader,
		walletWriter:  walletWriter,
	}
}

// IsExistingWallet validates if the Wallet exist
func (s *WalletDomainService) IsExistingWallet(walletID uuid.UUID) (bool, error) {
	exist, err := s.walletChecker.ExistByID(walletID)
	if err != nil {
		return false, err
	}
	return exist, nil
}

// MakePayment facilitates making a payment from source wallet
func (s *WalletDomainService) MakePayment(sourceWalletID uuid.UUID, destWalletID uuid.UUID, amount uint64) (*entity.Wallet, *entity.Payment, error) {

	sourceWallet, err := s.walletReader.GetByID(sourceWalletID)
	if err != nil {
		return nil, nil, err
	}
	if err := sourceWallet.Subtract(amount); err != nil {
		return nil, nil, err
	}

	payment, err := entity.NewPayment(sourceWalletID, destWalletID, amount, "OUT")
	if err != nil {
		return nil, nil, err
	}

	return sourceWallet, payment, nil
}

// ReceivePayment facilitates receiving a payment into source wallet
func (s *WalletDomainService) ReceivePayment(sourceWalletID uuid.UUID, destWalletID uuid.UUID, amount uint64) (*entity.Wallet, *entity.Payment, error) {

	sourceWallet, err := s.walletReader.GetByID(sourceWalletID)
	if err != nil {
		return nil, nil, err
	}
	if err := sourceWallet.Add(amount); err != nil {
		return nil, nil, err
	}

	payment, err := entity.NewPayment(sourceWalletID, destWalletID, amount, "IN")
	if err != nil {
		return nil, nil, err
	}

	return sourceWallet, payment, nil
}

// CreateWalletForNewUser creates a new wallet whenever a new UserAccount is created from Account module
// I might need to shift this later as it is abit odd for a domain service to be aware of Account module
func (s *WalletDomainService) CreateWalletForNewUser(sourceWalletID uuid.UUID) error {
	wallet, err := entity.NewWallet(sourceWalletID, 0)

	if err != nil {
		return err
	}

	if err := s.walletWriter.Create(wallet); err != nil {
		return err
	}
	return nil
}
