// internal/repository/transaction_repository.go
package repository

import (
	"GonPay_Backend/internal/domain"
	"database/sql"
)

type transactionRepository struct {
	db *PostgresDB
}

func NewTransactionRepository(db *PostgresDB) domain.TransactionRepository {
	return &transactionRepository{db: db}
}

func (r *transactionRepository) Create(tx *domain.Transaction) error {
	query := `
        INSERT INTO transactions 
        (source_wallet_id, destination_wallet_id, transaction_type, amount, reference_id, status, description)
        VALUES ($1, $2, $3, $4, uuid_generate_v4(), $5, $6)
        RETURNING transaction_id, reference_id, created_at`

	return r.db.DB.QueryRow(
		query,
		tx.SourceWalletID,
		tx.DestinationWalletID,
		tx.Type,
		tx.Amount,
		tx.Status,
		tx.Description,
	).Scan(&tx.ID, &tx.ReferenceID, &tx.CreatedAt)
}

func (r *transactionRepository) GetByID(id int64) (*domain.Transaction, error) {
	tx := &domain.Transaction{}
	query := `
        SELECT 
            transaction_id, source_wallet_id, destination_wallet_id, 
            transaction_type, amount, reference_id, status, 
            description, created_at
        FROM transactions 
        WHERE transaction_id = $1`

	err := r.db.DB.QueryRow(query, id).Scan(
		&tx.ID,
		&tx.SourceWalletID,
		&tx.DestinationWalletID,
		&tx.Type,
		&tx.Amount,
		&tx.ReferenceID,
		&tx.Status,
		&tx.Description,
		&tx.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, domain.ErrInvalidOperation
	}
	return tx, err
}

func (r *transactionRepository) GetByWalletID(walletID int64) ([]*domain.Transaction, error) {
	query := `
        SELECT 
            transaction_id, source_wallet_id, destination_wallet_id, 
            transaction_type, amount, reference_id, status, 
            description, created_at
        FROM transactions 
        WHERE source_wallet_id = $1 OR destination_wallet_id = $1
        ORDER BY created_at DESC`

	rows, err := r.db.DB.Query(query, walletID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []*domain.Transaction
	for rows.Next() {
		tx := &domain.Transaction{}
		err := rows.Scan(
			&tx.ID,
			&tx.SourceWalletID,
			&tx.DestinationWalletID,
			&tx.Type,
			&tx.Amount,
			&tx.ReferenceID,
			&tx.Status,
			&tx.Description,
			&tx.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, tx)
	}

	return transactions, nil
}

func (r *transactionRepository) UpdateStatus(id int64, status domain.TransactionStatus) error {
	query := `UPDATE transactions SET status = $1 WHERE transaction_id = $2`

	result, err := r.db.DB.Exec(query, status, id)
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

func (r *transactionRepository) GetUserTransactions(userID int64, limit, offset int) ([]*domain.Transaction, error) {
	query := `
        SELECT 
            t.transaction_id, t.source_wallet_id, t.destination_wallet_id,
            t.transaction_type, t.amount, t.reference_id, t.status,
            t.description, t.created_at
        FROM transactions t
        INNER JOIN wallets w ON t.source_wallet_id = w.wallet_id 
            OR t.destination_wallet_id = w.wallet_id
        WHERE w.user_id = $1 
        ORDER BY t.created_at DESC
        LIMIT $2 OFFSET $3`

	rows, err := r.db.DB.Query(query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []*domain.Transaction
	for rows.Next() {
		tx := &domain.Transaction{}
		var destWalletID sql.NullInt64 // Use sql.NullInt64 for nullable column

		err := rows.Scan(
			&tx.ID,
			&tx.SourceWalletID,
			&destWalletID,
			&tx.Type,
			&tx.Amount,
			&tx.ReferenceID,
			&tx.Status,
			&tx.Description,
			&tx.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		// Handle nullable destination_wallet_id
		if destWalletID.Valid {
			tx.DestinationWalletID = &destWalletID.Int64
		}

		transactions = append(transactions, tx)
	}

	return transactions, nil
}
