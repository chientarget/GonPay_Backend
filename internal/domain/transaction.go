// internal/domain/transaction.go
package domain

import (
	"time"
)

type TransactionType string
type TransactionStatus string

const (
	TransactionTypeDeposit  TransactionType = "DEPOSIT"
	TransactionTypeWithdraw TransactionType = "WITHDRAW"
	TransactionTypeTransfer TransactionType = "TRANSFER"

	TransactionStatusPending   TransactionStatus = "PENDING"
	TransactionStatusCompleted TransactionStatus = "COMPLETED"
	TransactionStatusFailed    TransactionStatus = "FAILED"
)

type Transaction struct {
	ID                  int64             `json:"id"`
	SourceWalletID      int64             `json:"source_wallet_id"`
	DestinationWalletID *int64            `json:"destination_wallet_id,omitempty"`
	Type                TransactionType   `json:"type"`
	Amount              float64           `json:"amount"`
	ReferenceID         string            `json:"reference_id"`
	Status              TransactionStatus `json:"status"`
	Description         string            `json:"description"`
	CreatedAt           time.Time         `json:"created_at"`
}

type TransactionRepository interface {
	Create(transaction *Transaction) error
	GetByID(id int64) (*Transaction, error)
	GetByWalletID(walletID int64) ([]*Transaction, error)
	UpdateStatus(id int64, status TransactionStatus) error
	GetUserTransactions(userID int64, limit, offset int) ([]*Transaction, error)
}
