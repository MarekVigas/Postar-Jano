package promo

import "errors"

var (
	ErrAlreadyUsed = errors.New("already used")
	ErrInvalid     = errors.New("invalid token")
)
