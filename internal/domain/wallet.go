// internal/domain/wallet.go
package domain

import (
	"time"
)

type Wallet struct {
	ID           int64      `json:"id"`
	UserID       int64      `json:"user_id"`
	WalletNumber string     `json:"wallet_number"`
	Balance      float64    `json:"balance"`
	Status       UserStatus `json:"status"`
	CreatedAt    time.Time  `json:"created_at"`
}

type WalletRepository interface {
	Create(wallet *Wallet) error
	GetByID(id int64) (*Wallet, error)
	GetByUserID(userID int64) ([]*Wallet, error)
	UpdateBalance(id int64, amount float64) error
	Delete(id int64) error
}
