// internal/domain/transaction_limit.go
package domain

import (
	"time"
)

type TransactionLimit struct {
	ID              int64           `json:"id"`
	UserID          int64           `json:"user_id"`
	TransactionType TransactionType `json:"transaction_type"`
	DailyLimit      float64         `json:"daily_limit"`
	MonthlyLimit    float64         `json:"monthly_limit"`
	CreatedAt       time.Time       `json:"created_at"`
}

type TransactionLimitRepository interface {
	Create(limit *TransactionLimit) error
	Update(limit *TransactionLimit) error
	GetByUserAndType(userID int64, transactionType TransactionType) (*TransactionLimit, error)
	GetByUserID(userID int64) ([]*TransactionLimit, error)
	CheckLimit(userID int64, transactionType TransactionType, amount float64) error
}
