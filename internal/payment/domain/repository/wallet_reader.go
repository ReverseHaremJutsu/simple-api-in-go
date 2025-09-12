package repository

import (
	"rest-api-in-gin/internal/payment/domain/model/entity"

	"github.com/google/uuid"
)

// WalletReader is a interface contract
type WalletReader interface {
	GetByID(id uuid.UUID) (*entity.Wallet, error)
}
