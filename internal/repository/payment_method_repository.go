// internal/repository/payment_method_repository.go
package repository

import (
	"GonPay_Backend/internal/domain"
	"database/sql"
)

type paymentMethodRepository struct {
	db *PostgresDB
}

func NewPaymentMethodRepository(db *PostgresDB) domain.PaymentMethodRepository {
	return &paymentMethodRepository{db: db}
}

func (r *paymentMethodRepository) Create(pm *domain.PaymentMethod) error {
	query := `
		INSERT INTO payment_methods (user_id, method_type, account_number, bank_name, is_default, status)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING payment_method_id, created_at`

	// Start a transaction if this is being set as default
	var err error
	if pm.IsDefault {
		// Update all other payment methods to non-default first
		updateQuery := `
			UPDATE payment_methods 
			SET is_default = false 
			WHERE user_id = $1 AND payment_method_id != COALESCE($2, 0)`

		if _, err = r.db.DB.Exec(updateQuery, pm.UserID, pm.ID); err != nil {
			return err
		}
	}

	return r.db.DB.QueryRow(
		query,
		pm.UserID,
		pm.MethodType,
		pm.AccountNumber,
		pm.BankName,
		pm.IsDefault,
		pm.Status,
	).Scan(&pm.ID, &pm.CreatedAt)
}

func (r *paymentMethodRepository) Update(pm *domain.PaymentMethod) error {
	tx, err := r.db.DB.Begin()
	if err != nil {
		return err
	}

	// If setting as default, update other payment methods first
	if pm.IsDefault {
		updateQuery := `
			UPDATE payment_methods 
			SET is_default = false 
			WHERE user_id = $1 AND payment_method_id != $2`

		if _, err = tx.Exec(updateQuery, pm.UserID, pm.ID); err != nil {
			tx.Rollback()
			return err
		}
	}

	query := `
		UPDATE payment_methods 
		SET method_type = $1, account_number = $2, bank_name = $3, is_default = $4, status = $5
		WHERE payment_method_id = $6 AND user_id = $7`

	result, err := tx.Exec(
		query,
		pm.MethodType,
		pm.AccountNumber,
		pm.BankName,
		pm.IsDefault,
		pm.Status,
		pm.ID,
		pm.UserID,
	)
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
		return domain.ErrInvalidOperation
	}

	return tx.Commit()
}

func (r *paymentMethodRepository) Delete(id int64) error {
	query := `UPDATE payment_methods SET status = $1 WHERE payment_method_id = $2`

	result, err := r.db.DB.Exec(query, domain.UserStatusInactive, id)
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

func (r *paymentMethodRepository) GetByID(id int64) (*domain.PaymentMethod, error) {
	pm := &domain.PaymentMethod{}
	query := `
		SELECT 
			payment_method_id, user_id, method_type, account_number, 
			bank_name, is_default, status, created_at
		FROM payment_methods 
		WHERE payment_method_id = $1`

	err := r.db.DB.QueryRow(query, id).Scan(
		&pm.ID,
		&pm.UserID,
		&pm.MethodType,
		&pm.AccountNumber,
		&pm.BankName,
		&pm.IsDefault,
		&pm.Status,
		&pm.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, domain.ErrInvalidOperation
	}
	return pm, err
}

func (r *paymentMethodRepository) GetByUserID(userID int64) ([]*domain.PaymentMethod, error) {
	query := `
		SELECT 
			payment_method_id, user_id, method_type, account_number, 
			bank_name, is_default, status, created_at
		FROM payment_methods 
		WHERE user_id = $1 AND status = $2
		ORDER BY is_default DESC, created_at DESC`

	rows, err := r.db.DB.Query(query, userID, domain.UserStatusActive)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var paymentMethods []*domain.PaymentMethod
	for rows.Next() {
		pm := &domain.PaymentMethod{}
		err := rows.Scan(
			&pm.ID,
			&pm.UserID,
			&pm.MethodType,
			&pm.AccountNumber,
			&pm.BankName,
			&pm.IsDefault,
			&pm.Status,
			&pm.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		paymentMethods = append(paymentMethods, pm)
	}

	return paymentMethods, nil
}

func (r *paymentMethodRepository) SetDefault(id int64, userID int64) error {
	tx, err := r.db.DB.Begin()
	if err != nil {
		return err
	}

	// First, set all payment methods for this user to non-default
	query1 := `
		UPDATE payment_methods 
		SET is_default = false 
		WHERE user_id = $1`

	if _, err = tx.Exec(query1, userID); err != nil {
		tx.Rollback()
		return err
	}

	// Then set the specified payment method as default
	query2 := `
		UPDATE payment_methods 
		SET is_default = true 
		WHERE payment_method_id = $1 AND user_id = $2`

	result, err := tx.Exec(query2, id, userID)
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
		return domain.ErrInvalidOperation
	}

	return tx.Commit()
}
