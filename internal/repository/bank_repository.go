package repository

import (
	"context"
	"database/sql"
	"errors"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInsufficientFunds  = errors.New("insufficient funds")
)

// BankRepository handles data access for bank domain.
type BankRepository struct {
	db *sql.DB
}

func NewBankRepository(db *sql.DB) *BankRepository {
	return &BankRepository{db: db}
}

func (r *BankRepository) Authenticate(ctx context.Context, email, password string) (int64, error) {
	var userID int64
	var storedPassword string
	query := "SELECT id, password_hash FROM users WHERE email = ?"
	err := r.db.QueryRowContext(ctx, query, email).Scan(&userID, &storedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrInvalidCredentials
		}
		return 0, err
	}

	if storedPassword != password {
		return 0, ErrInvalidCredentials
	}

	return userID, nil
}

func (r *BankRepository) FetchBalance(ctx context.Context, userID int64) (float64, error) {
	var balance float64
	query := "SELECT balance FROM accounts WHERE user_id = ?"
	if err := r.db.QueryRowContext(ctx, query, userID).Scan(&balance); err != nil {
		return 0, err
	}
	return balance, nil
}

func (r *BankRepository) CreateTransfer(ctx context.Context, fromUserID int64, toAccountNumber string, amount float64) (int64, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	var fromAccountID int64
	var balance float64
	query := "SELECT id, balance FROM accounts WHERE user_id = ? FOR UPDATE"
	if err := tx.QueryRowContext(ctx, query, fromUserID).Scan(&fromAccountID, &balance); err != nil {
		return 0, err
	}

	if balance < amount {
		return 0, ErrInsufficientFunds
	}

	update := "UPDATE accounts SET balance = balance - ? WHERE id = ?"
	if _, err := tx.ExecContext(ctx, update, amount, fromAccountID); err != nil {
		return 0, err
	}

	insert := "INSERT INTO transfers (from_account_id, to_account_number, amount) VALUES (?, ?, ?)"
	result, err := tx.ExecContext(ctx, insert, fromAccountID, toAccountNumber, amount)
	if err != nil {
		return 0, err
	}

	transferID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}

	return transferID, nil
}
