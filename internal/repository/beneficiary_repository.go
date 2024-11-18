// internal/repository/beneficiary_repository.go
package repository

import (
	"GonPay_Backend/internal/domain"
	"database/sql"
)

type beneficiaryRepository struct {
	db *PostgresDB
}

func NewBeneficiaryRepository(db *PostgresDB) domain.BeneficiaryRepository {
	return &beneficiaryRepository{db: db}
}

func (r *beneficiaryRepository) Create(b *domain.Beneficiary) error {
	query := `
        INSERT INTO beneficiaries (user_id, beneficiary_name, account_identifier, account_type, bank_name)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING beneficiary_id, created_at`

	return r.db.DB.QueryRow(
		query,
		b.UserID,
		b.BeneficiaryName,
		b.AccountIdentifier,
		b.AccountType,
		b.BankName,
	).Scan(&b.ID, &b.CreatedAt)
}

func (r *beneficiaryRepository) Update(b *domain.Beneficiary) error {
	query := `
        UPDATE beneficiaries 
        SET beneficiary_name = $1, account_identifier = $2, account_type = $3, bank_name = $4
        WHERE beneficiary_id = $5 AND user_id = $6`

	result, err := r.db.DB.Exec(
		query,
		b.BeneficiaryName,
		b.AccountIdentifier,
		b.AccountType,
		b.BankName,
		b.ID,
		b.UserID,
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

func (r *beneficiaryRepository) Delete(id int64) error {
	query := `DELETE FROM beneficiaries WHERE beneficiary_id = $1`

	result, err := r.db.DB.Exec(query, id)
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

func (r *beneficiaryRepository) GetByID(id int64) (*domain.Beneficiary, error) {
	b := &domain.Beneficiary{}
	query := `
        SELECT beneficiary_id, user_id, beneficiary_name, account_identifier, account_type, bank_name, created_at
        FROM beneficiaries 
        WHERE beneficiary_id = $1`

	err := r.db.DB.QueryRow(query, id).Scan(
		&b.ID,
		&b.UserID,
		&b.BeneficiaryName,
		&b.AccountIdentifier,
		&b.AccountType,
		&b.BankName,
		&b.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, domain.ErrInvalidOperation
	}
	return b, err
}

func (r *beneficiaryRepository) GetByUserID(userID int64) ([]*domain.Beneficiary, error) {
	query := `
        SELECT beneficiary_id, user_id, beneficiary_name, account_identifier, account_type, bank_name, created_at
        FROM beneficiaries 
        WHERE user_id = $1
        ORDER BY created_at DESC`

	rows, err := r.db.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var beneficiaries []*domain.Beneficiary
	for rows.Next() {
		b := &domain.Beneficiary{}
		err := rows.Scan(
			&b.ID,
			&b.UserID,
			&b.BeneficiaryName,
			&b.AccountIdentifier,
			&b.AccountType,
			&b.BankName,
			&b.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		beneficiaries = append(beneficiaries, b)
	}

	return beneficiaries, nil
}

func (r *beneficiaryRepository) GetByAccountIdentifier(accountIdentifier string, accountType domain.AccountType) (*domain.Beneficiary, error) {
	b := &domain.Beneficiary{}
	query := `
        SELECT beneficiary_id, user_id, beneficiary_name, account_identifier, account_type, bank_name, created_at
        FROM beneficiaries 
        WHERE account_identifier = $1 AND account_type = $2`

	err := r.db.DB.QueryRow(query, accountIdentifier, accountType).Scan(
		&b.ID,
		&b.UserID,
		&b.BeneficiaryName,
		&b.AccountIdentifier,
		&b.AccountType,
		&b.BankName,
		&b.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	return b, err
}
