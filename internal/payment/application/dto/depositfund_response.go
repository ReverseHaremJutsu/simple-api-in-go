package dto

import "github.com/google/uuid"

// DepositFundResponse is the output for the DepositFunds application service
type DepositFundResponse struct {
	WalletID       uuid.UUID `json:"wallet_id"`
	Amount         uint64    `json:"amount"`
	UpdatedBalance uint64    `json:"balance"`
}

// NewDepositFundResponse creates a new instance of DepositFundResponse
func NewDepositFundResponse(walletID uuid.UUID, amount uint64, updatedBalance uint64) *DepositFundResponse {
	return &DepositFundResponse{
		WalletID:       walletID,
		Amount:         amount,
		UpdatedBalance: updatedBalance,
	}
}
