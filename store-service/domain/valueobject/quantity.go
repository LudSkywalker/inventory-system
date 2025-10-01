package valueobject

import "errors"

var (
	ErrNegativeQuantity = errors.New("quantity cannot be negative")
)

type Quantity struct {
	value int
}

func NewQuantity(value int) (Quantity, error) {
	if value < 0 {
		return Quantity{}, ErrNegativeQuantity
	}
	return Quantity{value: value}, nil
}

func (q Quantity) Value() int {
	return q.value
}

func (q Quantity) IsValid() bool {
	return q.value >= 0
}

func (q Quantity) Add(other Quantity) (Quantity, error) {
	return NewQuantity(q.value + other.value)
}

func (q Quantity) Subtract(other Quantity) (Quantity, error) {
	return NewQuantity(q.value - other.value)
}
