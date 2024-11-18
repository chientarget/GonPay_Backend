// internal/domain/user.go
package domain

import (
	"time"
)

type UserStatus string

const (
	UserStatusActive   UserStatus = "ACTIVE"
	UserStatusInactive UserStatus = "INACTIVE"
)

const (
	RoleUser  = "USER"
	RoleAdmin = "ADMIN"
)

type User struct {
	ID           int64      `json:"id"`
	Username     string     `json:"username"`
	Email        string     `json:"email"`
	PhoneNumber  string     `json:"phone_number"`
	PasswordHash string     `json:"-"`
	Status       UserStatus `json:"status"`
	Preferences  string     `json:"preferences"`
	Role         string     `json:"role"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

type UserRepository interface {
	Create(user *User) error
	GetByID(id int64) (*User, error)
	GetByEmail(email string) (*User, error)
	Update(user *User) error
	Delete(id int64) error
}
