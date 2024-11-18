// internal/usecase/wallet_usecase.go
package usecase

import (
	"GonPay_Backend/internal/domain"
	"errors"
)

type WalletUseCase struct {
	walletRepo      domain.WalletRepository
	transactionRepo domain.TransactionRepository
}

func NewWalletUseCase(walletRepo domain.WalletRepository, transactionRepo domain.TransactionRepository) *WalletUseCase {
	return &WalletUseCase{
		walletRepo:      walletRepo,
		transactionRepo: transactionRepo,
	}
}

func (u *WalletUseCase) CreateWallet(userID int64) (*domain.Wallet, error) {
	wallet := &domain.Wallet{
		UserID:  userID,
		Balance: 0,
		Status:  domain.UserStatusActive,
	}

	if err := u.walletRepo.Create(wallet); err != nil {
		return nil, err
	}

	return wallet, nil
}

func (u *WalletUseCase) GetUserWallets(userID int64) ([]*domain.Wallet, error) {
	return u.walletRepo.GetByUserID(userID)
}

func (u *WalletUseCase) GetWallet(walletID int64) (*domain.Wallet, error) {
	return u.walletRepo.GetByID(walletID)
}

func (u *WalletUseCase) DeactivateWallet(walletID int64) error {
	wallet, err := u.walletRepo.GetByID(walletID)
	if err != nil {
		return err
	}

	if wallet.Balance > 0 {
		return errors.New("cannot deactivate wallet with positive balance")
	}

	return u.walletRepo.Delete(walletID)
}

func (u *WalletUseCase) Transfer(sourceWalletID, destWalletID int64, amount float64) (*domain.Transaction, error) {
	if amount <= 0 {
		return nil, domain.ErrInvalidAmount
	}

	// Create transaction record
	tx := &domain.Transaction{
		SourceWalletID:      sourceWalletID,
		DestinationWalletID: &destWalletID,
		Type:                domain.TransactionTypeTransfer,
		Amount:              amount,
		Status:              domain.TransactionStatusPending,
	}

	if err := u.transactionRepo.Create(tx); err != nil {
		return nil, err
	}

	// Update source wallet
	if err := u.walletRepo.UpdateBalance(sourceWalletID, -amount); err != nil {
		// Rollback transaction status if update fails
		u.transactionRepo.UpdateStatus(tx.ID, domain.TransactionStatusFailed)
		return nil, err
	}

	// Update destination wallet
	if err := u.walletRepo.UpdateBalance(destWalletID, amount); err != nil {
		// Rollback both transaction and source wallet if update fails
		u.walletRepo.UpdateBalance(sourceWalletID, amount)
		u.transactionRepo.UpdateStatus(tx.ID, domain.TransactionStatusFailed)
		return nil, err
	}

	// Update transaction status to completed
	if err := u.transactionRepo.UpdateStatus(tx.ID, domain.TransactionStatusCompleted); err != nil {
		return nil, err
	}

	return tx, nil
}

func (u *WalletUseCase) Deposit(walletID int64, amount float64) (*domain.Transaction, error) {
	if amount <= 0 {
		return nil, domain.ErrInvalidAmount
	}

	// Verify wallet exists
	if _, err := u.walletRepo.GetByID(walletID); err != nil {
		return nil, err
	}

	// Create transaction record
	tx := &domain.Transaction{
		SourceWalletID: walletID,
		Type:           domain.TransactionTypeDeposit,
		Amount:         amount,
		Status:         domain.TransactionStatusPending,
	}

	if err := u.transactionRepo.Create(tx); err != nil {
		return nil, err
	}

	// Update wallet balance
	if err := u.walletRepo.UpdateBalance(walletID, amount); err != nil {
		u.transactionRepo.UpdateStatus(tx.ID, domain.TransactionStatusFailed)
		return nil, err
	}

	// Update transaction status
	if err := u.transactionRepo.UpdateStatus(tx.ID, domain.TransactionStatusCompleted); err != nil {
		return nil, err
	}

	return tx, nil
}

func (u *WalletUseCase) Withdraw(walletID int64, amount float64) (*domain.Transaction, error) {
	if amount <= 0 {
		return nil, domain.ErrInvalidAmount
	}

	// Verify wallet exists and has sufficient funds
	wallet, err := u.walletRepo.GetByID(walletID)
	if err != nil {
		return nil, err
	}

	if wallet.Balance < amount {
		return nil, domain.ErrInsufficientFunds
	}

	// Create transaction record
	tx := &domain.Transaction{
		SourceWalletID: walletID,
		Type:           domain.TransactionTypeWithdraw,
		Amount:         amount,
		Status:         domain.TransactionStatusPending,
	}

	if err := u.transactionRepo.Create(tx); err != nil {
		return nil, err
	}

	// Update wallet balance
	if err := u.walletRepo.UpdateBalance(walletID, -amount); err != nil {
		u.transactionRepo.UpdateStatus(tx.ID, domain.TransactionStatusFailed)
		return nil, err
	}

	// Update transaction status
	if err := u.transactionRepo.UpdateStatus(tx.ID, domain.TransactionStatusCompleted); err != nil {
		return nil, err
	}

	return tx, nil
}
