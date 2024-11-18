// internal/usecase/payment_method_usecase.go
package usecase

import (
	"GonPay_Backend/internal/domain"
)

type PaymentMethodUseCase struct {
	paymentMethodRepo domain.PaymentMethodRepository
}

func NewPaymentMethodUseCase(paymentMethodRepo domain.PaymentMethodRepository) *PaymentMethodUseCase {
	return &PaymentMethodUseCase{
		paymentMethodRepo: paymentMethodRepo,
	}
}

func (u *PaymentMethodUseCase) CreatePaymentMethod(userID int64, methodType domain.PaymentMethodType, accountNumber string, bankName string, isDefault bool) (*domain.PaymentMethod, error) {
	pm := &domain.PaymentMethod{
		UserID:        userID,
		MethodType:    methodType,
		AccountNumber: accountNumber,
		BankName:      bankName,
		IsDefault:     isDefault,
		Status:        domain.UserStatusActive,
	}

	if err := u.paymentMethodRepo.Create(pm); err != nil {
		return nil, err
	}

	return pm, nil
}

func (u *PaymentMethodUseCase) UpdatePaymentMethod(id int64, userID int64, methodType domain.PaymentMethodType, accountNumber string, bankName string, isDefault bool) (*domain.PaymentMethod, error) {
	pm, err := u.paymentMethodRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if pm.UserID != userID {
		return nil, domain.ErrInvalidOperation
	}

	pm.MethodType = methodType
	pm.AccountNumber = accountNumber
	pm.BankName = bankName
	pm.IsDefault = isDefault

	if err := u.paymentMethodRepo.Update(pm); err != nil {
		return nil, err
	}

	return pm, nil
}

func (u *PaymentMethodUseCase) DeletePaymentMethod(id int64, userID int64) error {
	pm, err := u.paymentMethodRepo.GetByID(id)
	if err != nil {
		return err
	}

	if pm.UserID != userID {
		return domain.ErrInvalidOperation
	}

	return u.paymentMethodRepo.Delete(id)
}

func (u *PaymentMethodUseCase) GetPaymentMethod(id int64, userID int64) (*domain.PaymentMethod, error) {
	pm, err := u.paymentMethodRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if pm.UserID != userID {
		return nil, domain.ErrInvalidOperation
	}

	return pm, nil
}

func (u *PaymentMethodUseCase) GetUserPaymentMethods(userID int64) ([]*domain.PaymentMethod, error) {
	return u.paymentMethodRepo.GetByUserID(userID)
}

func (u *PaymentMethodUseCase) SetDefaultPaymentMethod(id int64, userID int64) error {
	pm, err := u.paymentMethodRepo.GetByID(id)
	if err != nil {
		return err
	}

	if pm.UserID != userID {
		return domain.ErrInvalidOperation
	}

	return u.paymentMethodRepo.SetDefault(id, userID)
}
