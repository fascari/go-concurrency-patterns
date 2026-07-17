package mutex

import "errors"

var (
	ErrInvalidAmount     = errors.New("transfer amount must be positive")
	ErrSelfTransfer      = errors.New("source and destination accounts must differ")
	ErrAccountNotFound   = errors.New("account not found")
	ErrInsufficientFunds = errors.New("insufficient funds")
)
