// internal/usecase/transaction_usecase.go
package usecase

import (
	"GonPay_Backend/internal/domain"
	"errors"
)

type TransactionUseCase struct {
	transactionRepo domain.TransactionRepository
}

func NewTransactionUseCase(transactionRepo domain.TransactionRepository) *TransactionUseCase {
	return &TransactionUseCase{
		transactionRepo: transactionRepo,
	}
}

func (u *TransactionUseCase) GetUserTransactions(userID int64, page, limit int) ([]*domain.Transaction, error) {
	offset := (page - 1) * limit

	// Add validation
	if page < 1 {
		return nil, errors.New("page must be greater than 0")
	}

	if limit < 1 || limit > 100 {
		limit = 10 // Default limit
	}

	transactions, err := u.transactionRepo.GetUserTransactions(userID, limit, offset)
	if err != nil {
		return nil, err
	}

	// Return empty array instead of null if no transactions
	if transactions == nil {
		return []*domain.Transaction{}, nil
	}

	return transactions, nil
}
