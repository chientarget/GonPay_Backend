// internal/usecase/transaction_limit_usecase.go
package usecase

import (
	"GonPay_Backend/internal/domain"
	"errors"
)

type TransactionLimitUseCase struct {
	limitRepo domain.TransactionLimitRepository
}

func NewTransactionLimitUseCase(limitRepo domain.TransactionLimitRepository) *TransactionLimitUseCase {
	return &TransactionLimitUseCase{
		limitRepo: limitRepo,
	}
}

func (u *TransactionLimitUseCase) SetTransactionLimit(userID int64, transactionType domain.TransactionType, dailyLimit, monthlyLimit float64) (*domain.TransactionLimit, error) {
	// Validate limits
	if dailyLimit <= 0 || monthlyLimit <= 0 {
		return nil, errors.New("limits must be greater than 0")
	}

	if dailyLimit > monthlyLimit {
		return nil, errors.New("daily limit cannot exceed monthly limit")
	}

	// Check if limit exists
	existingLimit, err := u.limitRepo.GetByUserAndType(userID, transactionType)
	if err != nil {
		return nil, err
	}

	limit := &domain.TransactionLimit{
		UserID:          userID,
		TransactionType: transactionType,
		DailyLimit:      dailyLimit,
		MonthlyLimit:    monthlyLimit,
	}

	if existingLimit != nil {
		limit.ID = existingLimit.ID
		err = u.limitRepo.Update(limit)
	} else {
		err = u.limitRepo.Create(limit)
	}

	if err != nil {
		return nil, err
	}

	return limit, nil
}

func (u *TransactionLimitUseCase) GetUserLimits(userID int64) ([]*domain.TransactionLimit, error) {
	return u.limitRepo.GetByUserID(userID)
}

func (u *TransactionLimitUseCase) GetLimitByType(userID int64, transactionType domain.TransactionType) (*domain.TransactionLimit, error) {
	return u.limitRepo.GetByUserAndType(userID, transactionType)
}

func (u *TransactionLimitUseCase) CheckTransactionLimit(userID int64, transactionType domain.TransactionType, amount float64) error {
	if amount <= 0 {
		return domain.ErrInvalidAmount
	}

	return u.limitRepo.CheckLimit(userID, transactionType, amount)
}
