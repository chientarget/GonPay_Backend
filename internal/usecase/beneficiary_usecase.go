// internal/usecase/beneficiary_usecase.go
package usecase

import (
	"GonPay_Backend/internal/domain"
	"errors"
)

type BeneficiaryUseCase struct {
	beneficiaryRepo domain.BeneficiaryRepository
}

func NewBeneficiaryUseCase(beneficiaryRepo domain.BeneficiaryRepository) *BeneficiaryUseCase {
	return &BeneficiaryUseCase{
		beneficiaryRepo: beneficiaryRepo,
	}
}

func (u *BeneficiaryUseCase) CreateBeneficiary(userID int64, name string, accountIdentifier string, accountType domain.AccountType, bankName string) (*domain.Beneficiary, error) {
	// Check if beneficiary already exists
	existing, err := u.beneficiaryRepo.GetByAccountIdentifier(accountIdentifier, accountType)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("beneficiary with this account already exists")
	}

	beneficiary := &domain.Beneficiary{
		UserID:            userID,
		BeneficiaryName:   name,
		AccountIdentifier: accountIdentifier,
		AccountType:       accountType,
		BankName:          bankName,
	}

	if err := u.beneficiaryRepo.Create(beneficiary); err != nil {
		return nil, err
	}

	return beneficiary, nil
}

func (u *BeneficiaryUseCase) UpdateBeneficiary(id int64, userID int64, name string, accountIdentifier string, accountType domain.AccountType, bankName string) (*domain.Beneficiary, error) {
	// Get existing beneficiary
	beneficiary, err := u.beneficiaryRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Check ownership
	if beneficiary.UserID != userID {
		return nil, domain.ErrInvalidOperation
	}

	// Check if new account identifier already exists (if changed)
	if beneficiary.AccountIdentifier != accountIdentifier || beneficiary.AccountType != accountType {
		existing, err := u.beneficiaryRepo.GetByAccountIdentifier(accountIdentifier, accountType)
		if err != nil {
			return nil, err
		}
		if existing != nil && existing.ID != id {
			return nil, errors.New("beneficiary with this account already exists")
		}
	}

	beneficiary.BeneficiaryName = name
	beneficiary.AccountIdentifier = accountIdentifier
	beneficiary.AccountType = accountType
	beneficiary.BankName = bankName

	if err := u.beneficiaryRepo.Update(beneficiary); err != nil {
		return nil, err
	}

	return beneficiary, nil
}

func (u *BeneficiaryUseCase) DeleteBeneficiary(id int64, userID int64) error {
	// Get existing beneficiary
	beneficiary, err := u.beneficiaryRepo.GetByID(id)
	if err != nil {
		return err
	}

	// Check ownership
	if beneficiary.UserID != userID {
		return domain.ErrInvalidOperation
	}

	return u.beneficiaryRepo.Delete(id)
}

func (u *BeneficiaryUseCase) GetBeneficiary(id int64, userID int64) (*domain.Beneficiary, error) {
	beneficiary, err := u.beneficiaryRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if beneficiary.UserID != userID {
		return nil, domain.ErrInvalidOperation
	}

	return beneficiary, nil
}

func (u *BeneficiaryUseCase) GetUserBeneficiaries(userID int64) ([]*domain.Beneficiary, error) {
	return u.beneficiaryRepo.GetByUserID(userID)
}
