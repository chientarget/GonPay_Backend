// internal/domain/errors.go
package domain

import "errors"

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserNotFound       = errors.New("user not found")
	ErrWalletNotFound     = errors.New("wallet not found")
	ErrInsufficientFunds  = errors.New("insufficient funds")
	ErrInvalidAmount      = errors.New("invalid amount")
	ErrInvalidOperation   = errors.New("invalid operation")
)
