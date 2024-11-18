// internal/repository/transaction_limit_repository.go
package repository

import (
	"GonPay_Backend/internal/domain"
	"database/sql"
	_ "encoding/json"
	_ "errors"
	_ "net/http"
	_ "time"
)

type transactionLimitRepository struct {
	db *PostgresDB
}

func NewTransactionLimitRepository(db *PostgresDB) domain.TransactionLimitRepository {
	return &transactionLimitRepository{db: db}
}

func (r *transactionLimitRepository) Create(limit *domain.TransactionLimit) error {
	query := `
        INSERT INTO transaction_limits (user_id, transaction_type, daily_limit, monthly_limit)
        VALUES ($1, $2, $3, $4)
        RETURNING limit_id, created_at`

	return r.db.DB.QueryRow(
		query,
		limit.UserID,
		limit.TransactionType,
		limit.DailyLimit,
		limit.MonthlyLimit,
	).Scan(&limit.ID, &limit.CreatedAt)
}

func (r *transactionLimitRepository) Update(limit *domain.TransactionLimit) error {
	query := `
        UPDATE transaction_limits 
        SET daily_limit = $1, monthly_limit = $2
        WHERE limit_id = $3 AND user_id = $4`

	result, err := r.db.DB.Exec(
		query,
		limit.DailyLimit,
		limit.MonthlyLimit,
		limit.ID,
		limit.UserID,
	)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return domain.ErrInvalidOperation
	}

	return nil
}

func (r *transactionLimitRepository) GetByUserAndType(userID int64, transactionType domain.TransactionType) (*domain.TransactionLimit, error) {
	limit := &domain.TransactionLimit{}
	query := `
        SELECT limit_id, user_id, transaction_type, daily_limit, monthly_limit, created_at
        FROM transaction_limits 
        WHERE user_id = $1 AND transaction_type = $2`

	err := r.db.DB.QueryRow(query, userID, transactionType).Scan(
		&limit.ID,
		&limit.UserID,
		&limit.TransactionType,
		&limit.DailyLimit,
		&limit.MonthlyLimit,
		&limit.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	return limit, err
}

func (r *transactionLimitRepository) GetByUserID(userID int64) ([]*domain.TransactionLimit, error) {
	query := `
        SELECT limit_id, user_id, transaction_type, daily_limit, monthly_limit, created_at
        FROM transaction_limits 
        WHERE user_id = $1
        ORDER BY created_at DESC`

	rows, err := r.db.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var limits []*domain.TransactionLimit
	for rows.Next() {
		limit := &domain.TransactionLimit{}
		err := rows.Scan(
			&limit.ID,
			&limit.UserID,
			&limit.TransactionType,
			&limit.DailyLimit,
			&limit.MonthlyLimit,
			&limit.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		limits = append(limits, limit)
	}

	return limits, nil
}

func (r *transactionLimitRepository) CheckLimit(userID int64, transactionType domain.TransactionType, amount float64) error {
	// Check daily limit
	dailyQuery := `
        SELECT COALESCE(SUM(amount), 0)
        FROM transactions
        WHERE source_wallet_id IN (SELECT wallet_id FROM wallets WHERE user_id = $1)
        AND transaction_type = $2
        AND created_at >= CURRENT_DATE`

	var dailyTotal float64
	if err := r.db.DB.QueryRow(dailyQuery, userID, transactionType).Scan(&dailyTotal); err != nil {
		return err
	}

	// Check monthly limit
	monthlyQuery := `
        SELECT COALESCE(SUM(amount), 0)
        FROM transactions
        WHERE source_wallet_id IN (SELECT wallet_id FROM wallets WHERE user_id = $1)
        AND transaction_type = $2
        AND created_at >= DATE_TRUNC('month', CURRENT_DATE)`

	var monthlyTotal float64
	if err := r.db.DB.QueryRow(monthlyQuery, userID, transactionType).Scan(&monthlyTotal); err != nil {
		return err
	}

	// Get limits
	limit, err := r.GetByUserAndType(userID, transactionType)
	if err != nil {
		return err
	}

	// If no limits set, allow transaction
	if limit == nil {
		return nil
	}

	// Check if transaction would exceed limits
	if dailyTotal+amount > limit.DailyLimit {
		return domain.ErrInvalidOperation
	}

	if monthlyTotal+amount > limit.MonthlyLimit {
		return domain.ErrInvalidOperation
	}

	return nil
}
