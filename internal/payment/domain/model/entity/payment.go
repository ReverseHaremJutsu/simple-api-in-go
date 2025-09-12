package entity

import (
	"fmt"
	"rest-api-in-gin/internal/payment/domain/model/value_object"

	"github.com/google/uuid"
)

type PaymentAction string

const (
	PaymentIn  PaymentAction = "IN"
	PaymentOut PaymentAction = "OUT"
)

// Payment is an Entity
type Payment struct {
	ID                  uuid.UUID
	sourceWalletID      uuid.UUID
	destinationWalletID uuid.UUID
	amount              value_object.Money
	paymentAction       PaymentAction
}

// NewPayment validates and returns a new Payment Entity
func NewPayment(source uuid.UUID, destination uuid.UUID, amount uint64, action string) (*Payment, error) {
	moneyVO, err := value_object.NewMoney(amount)
	if err != nil {
		return nil, fmt.Errorf("error creating payment: %w", err)
	}

	paymentActionConst, err := NewPaymentAction(action)
	if err != nil {
		return nil, fmt.Errorf("error creating payment: %w", err)
	}

	return &Payment{
		ID:                  uuid.New(),
		sourceWalletID:      source,
		destinationWalletID: destination,
		amount:              moneyVO,
		paymentAction:       paymentActionConst,
	}, nil

}

// SourceWallet returns the uuid of the source wallet of Payment
func (p Payment) SourceWallet() uuid.UUID {
	return p.sourceWalletID
}

// DestinationWallet returns the uuid of the destination wallet of Paymemt
func (p Payment) DestinationWallet() uuid.UUID {
	return p.destinationWalletID
}

// Amount returns amount of Payment VO
func (p Payment) Amount() uint64 {
	return p.amount.Value()
}

// PaymentAction returns the string value of paymentAction
func (p Payment) PaymentAction() string {
	return string(p.paymentAction)
}

// NewPaymentAction creates a PaymentAction
func NewPaymentAction(action string) (PaymentAction, error) {
	switch action {
	case string(PaymentIn), string(PaymentOut):
		return PaymentAction(action), nil
	default:
		return "", fmt.Errorf("invalid payment action: %s", action)
	}
}

// EqualTo compares two Payment entities
func (p Payment) EqualTo(other Payment) bool {
	return p.ID == other.ID
}
