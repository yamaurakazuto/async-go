package service

import (
	"context"
	"errors"

	"async-go/internal/repository"
)

var ErrInsufficientBalance = errors.New("insufficient balance")

// BankService handles business logic.
type BankService struct {
	repo *repository.BankRepository
}

func NewBankService(repo *repository.BankRepository) *BankService {
	return &BankService{repo: repo}
}

func (s *BankService) Login(ctx context.Context, email, password string) (int64, error) {
	return s.repo.Authenticate(ctx, email, password)
}

func (s *BankService) Balance(ctx context.Context, userID int64) (float64, error) {
	return s.repo.FetchBalance(ctx, userID)
}

func (s *BankService) Transfer(ctx context.Context, fromUserID int64, toAccountNumber string, amount float64) (int64, error) {
	if amount <= 0 {
		return 0, errors.New("amount must be positive")
	}
	transferID, err := s.repo.CreateTransfer(ctx, fromUserID, toAccountNumber, amount)
	if err != nil {
		if errors.Is(err, repository.ErrInsufficientFunds) {
			return 0, ErrInsufficientBalance
		}
		return 0, err
	}
	return transferID, nil
}
