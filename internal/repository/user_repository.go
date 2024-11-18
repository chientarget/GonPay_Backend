// internal/repository/user_repository.go
package repository

import (
	"GonPay_Backend/internal/domain"
	"database/sql"
)

type userRepository struct {
	db *PostgresDB
}

func NewUserRepository(db *PostgresDB) domain.UserRepository {
	return &userRepository{db: db}
}

//func (r *userRepository) Create(user *domain.User) error {
//	query := `
//		INSERT INTO users (username, email, phone_number, password_hash, status, preferences)
//		VALUES ($1, $2, $3, $4, $5, $6)
//		RETURNING user_id, created_at, updated_at`
//
//	return r.db.DB.QueryRow(
//		query,
//		user.Username,
//		user.Email,
//		user.PhoneNumber,
//		user.PasswordHash,
//		user.Status,
//		user.Preferences,
//	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
//}

func (r *userRepository) Create(user *domain.User) error {
	query := `
        INSERT INTO users (username, email, phone_number, password_hash, status, preferences, role)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
        RETURNING user_id, created_at, updated_at`

	return r.db.DB.QueryRow(
		query,
		user.Username,
		user.Email,
		user.PhoneNumber,
		user.PasswordHash,
		user.Status,
		user.Preferences,
		user.Role,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
}

func (r *userRepository) GetByID(id int64) (*domain.User, error) {
	user := &domain.User{}
	query := `
        SELECT user_id, username, email, phone_number, password_hash, status, preferences, role, created_at, updated_at
        FROM users 
        WHERE user_id = $1`

	err := r.db.DB.QueryRow(query, id).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PhoneNumber,
		&user.PasswordHash,
		&user.Status,
		&user.Preferences,
		&user.Role, // Thêm role vào đây
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, domain.ErrUserNotFound
	}
	return user, err
}

func (r *userRepository) GetByEmail(email string) (*domain.User, error) {
	user := &domain.User{}
	query := `
        SELECT user_id, username, email, phone_number, password_hash, status, preferences, role, created_at, updated_at
        FROM users 
        WHERE email = $1`

	err := r.db.DB.QueryRow(query, email).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PhoneNumber,
		&user.PasswordHash,
		&user.Status,
		&user.Preferences,
		&user.Role, // Thêm role vào đây
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, domain.ErrUserNotFound
	}
	return user, err
}

func (r *userRepository) Update(user *domain.User) error {
	query := `
        UPDATE users 
        SET username = $1, email = $2, phone_number = $3, status = $4, preferences = $5, role = $6, updated_at = $7
        WHERE user_id = $8`

	result, err := r.db.DB.Exec(
		query,
		user.Username,
		user.Email,
		user.PhoneNumber,
		user.Status,
		user.Preferences,
		user.Role, // Thêm role vào đây
		user.UpdatedAt,
		user.ID,
	)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return domain.ErrUserNotFound
	}

	return nil
}

func (r *userRepository) Delete(id int64) error {
	query := `UPDATE users SET status = $1, updated_at = CURRENT_TIMESTAMP WHERE user_id = $2`

	result, err := r.db.DB.Exec(query, domain.UserStatusInactive, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return domain.ErrUserNotFound
	}

	return nil
}
