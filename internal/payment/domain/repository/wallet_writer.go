package repository

import (
	"rest-api-in-gin/internal/payment/domain/model/entity"
)

type WalletPaymentUpdate struct {
	Wallet  *entity.Wallet
	Payment *entity.Payment
}

// PaymentWriter is a interface contract
type WalletWriter interface {
	UpdatesByID(updatePayments []WalletPaymentUpdate) error
	Create(wallet *entity.Wallet) error
}
