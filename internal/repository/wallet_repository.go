// internal/repository/wallet_repository.go
package repository

import (
	"GonPay_Backend/internal/domain"
	"database/sql"
)

type walletRepository struct {
	db *PostgresDB
}

func NewWalletRepository(db *PostgresDB) domain.WalletRepository {
	return &walletRepository{db: db}
}

func (r *walletRepository) Create(wallet *domain.Wallet) error {
	query := `
        INSERT INTO wallets (user_id, status, balance)
        VALUES ($1, $2, $3)
        RETURNING wallet_id, wallet_number, created_at`

	return r.db.DB.QueryRow(
		query,
		wallet.UserID,
		wallet.Status,
		wallet.Balance,
	).Scan(&wallet.ID, &wallet.WalletNumber, &wallet.CreatedAt)
}

func (r *walletRepository) GetByID(id int64) (*domain.Wallet, error) {
	wallet := &domain.Wallet{}
	query := `
        SELECT wallet_id, user_id, wallet_number, balance, status, created_at
        FROM wallets 
        WHERE wallet_id = $1`

	err := r.db.DB.QueryRow(query, id).Scan(
		&wallet.ID,
		&wallet.UserID,
		&wallet.WalletNumber,
		&wallet.Balance,
		&wallet.Status,
		&wallet.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, domain.ErrWalletNotFound
	}
	return wallet, err
}

func (r *walletRepository) GetByUserID(userID int64) ([]*domain.Wallet, error) {
	query := `
        SELECT wallet_id, user_id, wallet_number, balance, status, created_at
        FROM wallets 
        WHERE user_id = $1
        ORDER BY created_at DESC`

	rows, err := r.db.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var wallets []*domain.Wallet
	for rows.Next() {
		wallet := &domain.Wallet{}
		err := rows.Scan(
			&wallet.ID,
			&wallet.UserID,
			&wallet.WalletNumber,
			&wallet.Balance,
			&wallet.Status,
			&wallet.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		wallets = append(wallets, wallet)
	}

	return wallets, nil
}

func (r *walletRepository) UpdateBalance(id int64, amount float64) error {
	tx, err := r.db.DB.Begin()
	if err != nil {
		return err
	}

	// Get current balance with lock
	var currentBalance float64
	query := `
        SELECT balance 
        FROM wallets 
        WHERE wallet_id = $1 AND status = 'ACTIVE'
        FOR UPDATE`

	err = tx.QueryRow(query, id).Scan(&currentBalance)
	if err != nil {
		tx.Rollback()
		if err == sql.ErrNoRows {
			return domain.ErrWalletNotFound
		}
		return err
	}

	// Check if new balance would be negative
	if currentBalance+amount < 0 {
		tx.Rollback()
		return domain.ErrInsufficientFunds
	}

	// Update balance
	query = `
        UPDATE wallets 
        SET balance = balance + $1
        WHERE wallet_id = $2 AND status = 'ACTIVE'`

	result, err := tx.Exec(query, amount, id)
	if err != nil {
		tx.Rollback()
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		tx.Rollback()
		return err
	}

	if rows == 0 {
		tx.Rollback()
		return domain.ErrWalletNotFound
	}

	return tx.Commit()
}

func (r *walletRepository) Delete(id int64) error {
	query := `UPDATE wallets SET status = $1 WHERE wallet_id = $2`

	result, err := r.db.DB.Exec(query, domain.UserStatusInactive, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return domain.ErrWalletNotFound
	}

	return nil
}
