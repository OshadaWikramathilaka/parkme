package services

import (
	"context"
	"errors"
	"time"

	"github.com/dfanso/parkme-backend/internal/models"
	"github.com/dfanso/parkme-backend/internal/repositories"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	ErrInsufficientBalance = errors.New("insufficient balance")
	ErrWalletNotFound      = errors.New("wallet not found")
	ErrInvalidAmount       = errors.New("invalid amount")
)

type WalletService struct {
	repo *repositories.WalletRepository
}

func NewWalletService(repo *repositories.WalletRepository) *WalletService {
	return &WalletService{
		repo: repo,
	}
}

// GetOrCreateWallet gets the user's wallet or creates one if it doesn't exist
func (s *WalletService) GetOrCreateWallet(ctx context.Context, userID primitive.ObjectID) (*models.Wallet, error) {
	wallet, err := s.repo.FindByUserID(ctx, userID)
	if err == nil {
		return wallet, nil
	}

	if err != repositories.ErrNotFound {
		return nil, err
	}

	// Create new wallet
	wallet = &models.Wallet{
		UserID:    userID,
		Balance:   0,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.repo.Create(ctx, wallet); err != nil {
		return nil, err
	}

	return wallet, nil
}

// TopUp adds points to the wallet
func (s *WalletService) TopUp(ctx context.Context, userID primitive.ObjectID, amount float64) (*models.Transaction, error) {
	if amount <= 0 {
		return nil, ErrInvalidAmount
	}

	wallet, err := s.GetOrCreateWallet(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Create transaction
	transaction := &models.Transaction{
		WalletID:    wallet.ID,
		Type:        models.TransactionTypeTopUp,
		Amount:      amount,
		Description: "Wallet top-up",
		CreatedAt:   time.Now(),
	}

	if err := s.repo.CreateTransaction(ctx, transaction); err != nil {
		return nil, err
	}

	// Update wallet balance
	wallet.Balance += amount
	wallet.UpdatedAt = time.Now()
	if err := s.repo.Update(ctx, wallet); err != nil {
		return nil, err
	}

	return transaction, nil
}

// Deduct removes points from the wallet
func (s *WalletService) Deduct(ctx context.Context, userID primitive.ObjectID, amount float64, description string) (*models.Transaction, error) {
	if amount <= 0 {
		return nil, ErrInvalidAmount
	}

	wallet, err := s.repo.FindByUserID(ctx, userID)
	if err != nil {
		if err == repositories.ErrNotFound {
			return nil, ErrWalletNotFound
		}
		return nil, err
	}

	if wallet.Balance < amount {
		return nil, ErrInsufficientBalance
	}

	// Create transaction
	transaction := &models.Transaction{
		WalletID:    wallet.ID,
		Type:        models.TransactionTypeDeduct,
		Amount:      amount,
		Description: description,
		CreatedAt:   time.Now(),
	}

	if err := s.repo.CreateTransaction(ctx, transaction); err != nil {
		return nil, err
	}

	// Update wallet balance
	wallet.Balance -= amount
	wallet.UpdatedAt = time.Now()
	if err := s.repo.Update(ctx, wallet); err != nil {
		return nil, err
	}

	return transaction, nil
}

// GetTransactions returns all transactions for a wallet
func (s *WalletService) GetTransactions(ctx context.Context, userID primitive.ObjectID) ([]models.Transaction, error) {
	wallet, err := s.GetOrCreateWallet(ctx, userID)
	if err != nil {
		return nil, err
	}

	return s.repo.GetTransactions(ctx, wallet.ID)
}

// GetBalance returns the current wallet balance
func (s *WalletService) GetBalance(ctx context.Context, userID primitive.ObjectID) (float64, error) {
	wallet, err := s.GetOrCreateWallet(ctx, userID)
	if err != nil {
		return 0, err
	}

	return wallet.Balance, nil
}
