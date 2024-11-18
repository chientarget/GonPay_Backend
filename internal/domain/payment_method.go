// internal/domain/payment_method.go
package domain

import (
	"time"
)

type PaymentMethodType string

const (
	PaymentMethodTypeCreditCard  PaymentMethodType = "CREDIT_CARD"
	PaymentMethodTypeDebitCard   PaymentMethodType = "DEBIT_CARD"
	PaymentMethodTypeEWallet     PaymentMethodType = "E_WALLET"
	PaymentMethodTypeBankAccount PaymentMethodType = "BANK_ACCOUNT"
)

type PaymentMethod struct {
	ID            int64             `json:"id"`
	UserID        int64             `json:"user_id"`
	MethodType    PaymentMethodType `json:"method_type"`
	AccountNumber string            `json:"account_number"`
	BankName      string            `json:"bank_name,omitempty"`
	IsDefault     bool              `json:"is_default"`
	Status        UserStatus        `json:"status"`
	CreatedAt     time.Time         `json:"created_at"`
}

type PaymentMethodRepository interface {
	Create(paymentMethod *PaymentMethod) error
	Update(paymentMethod *PaymentMethod) error
	Delete(id int64) error
	GetByID(id int64) (*PaymentMethod, error)
	GetByUserID(userID int64) ([]*PaymentMethod, error)
	SetDefault(id int64, userID int64) error
}
