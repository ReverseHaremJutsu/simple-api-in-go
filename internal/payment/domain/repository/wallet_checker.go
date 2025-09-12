package repository

import "github.com/google/uuid"

// PaymentWriter is a interface contract
type WalletChecker interface {
	ExistByID(walletID uuid.UUID) (bool, error)
}
