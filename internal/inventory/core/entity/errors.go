package entity

import "errors"

var (
	ErrInvalidInput    = errors.New("invalid input: item ID and store ID are required")
	ErrInvalidQuantity = errors.New("invalid quantity value")
)
