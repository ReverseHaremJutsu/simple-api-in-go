package entity

import (
	"errors"
	"fmt"
	"rest-api-in-gin/internal/payment/domain/model/value_object"

	"github.com/google/uuid"
)

// Wallet is an Aggregate Entity
type Wallet struct {
	ID            uuid.UUID
	UserAccountID uuid.UUID
	Balance       value_object.Money
}

// NewWallet creates a new Wallet
func NewWallet(userAccountID uuid.UUID, balance uint64) (*Wallet, error) {
	moneyVO, err := value_object.NewMoney(balance)
	if err != nil {
		return nil, fmt.Errorf("error creating account: %w", err)
	}

	return &Wallet{
		ID:            uuid.New(),
		UserAccountID: userAccountID,
		Balance:       moneyVO,
	}, nil
}

// Subract decreases and updates the Wallet balance
func (w *Wallet) Subtract(amount uint64) error {
	if w.Balance.Value() < amount {
		return errors.New("insufficient funds")
	}
	newBalance, err := value_object.NewMoney(w.Balance.Value() - amount)
	if err != nil {
		return fmt.Errorf("error subtracting amount from wallet: %w", err)
	}

	w.Balance = newBalance
	return nil
}

// Add increases and updates the Wallet balance
func (w *Wallet) Add(amount uint64) error {
	newBalance, err := value_object.NewMoney(w.Balance.Value() + amount)
	if err != nil {
		return fmt.Errorf("error adding amount into wallet: %w", err)
	}

	w.Balance = newBalance
	return nil
}

// EqualTo
func (w Wallet) EqualTo(other Wallet) bool {
	return w.ID == other.ID
}
