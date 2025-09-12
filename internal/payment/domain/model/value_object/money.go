package value_object

// Money is a Value Object
type Money struct {
	amount uint64
}

// NewMoney returns a new Money VO
func NewMoney(amount uint64) (Money, error) {
	return Money{amount: amount}, nil
}

// Value returns amount of Money
func (m Money) Value() uint64 {
	return m.amount
}

// EqualTo compares two Money VO
func (m Money) EqualTo(other Money) bool {
	return m.amount == other.amount
}
