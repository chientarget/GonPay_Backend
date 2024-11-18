// internal/domain/beneficiary.go
package domain

import (
	"time"
)

type AccountType string

const (
	AccountTypeWallet      AccountType = "WALLET"
	AccountTypeBankAccount AccountType = "BANK_ACCOUNT"
)

type Beneficiary struct {
	ID                int64       `json:"id"`
	UserID            int64       `json:"user_id"`
	BeneficiaryName   string      `json:"beneficiary_name"`
	AccountIdentifier string      `json:"account_identifier"`
	AccountType       AccountType `json:"account_type"`
	BankName          string      `json:"bank_name,omitempty"`
	CreatedAt         time.Time   `json:"created_at"`
}

type BeneficiaryRepository interface {
	Create(beneficiary *Beneficiary) error
	Update(beneficiary *Beneficiary) error
	Delete(id int64) error
	GetByID(id int64) (*Beneficiary, error)
	GetByUserID(userID int64) ([]*Beneficiary, error)
	GetByAccountIdentifier(accountIdentifier string, accountType AccountType) (*Beneficiary, error)
}
