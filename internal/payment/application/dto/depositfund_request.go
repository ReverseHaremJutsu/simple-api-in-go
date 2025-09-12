package dto

import "github.com/google/uuid"

// DepositFundRequest is the input for the DepositFunds application service
type DepositFundRequest struct {
	WalletID uuid.UUID `json:"wallet_id"`
	Amount   uint64    `json:"amount"`
}
