package data

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type Wallet struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	Balance   float64   `json:"balance"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type WalletModel struct {
	DB *sql.DB
}

func (m WalletModel) Insert(wallet *Wallet) error {
	query := `
        INSERT INTO wallets (user_id, balance, created_at, updated_at)
        VALUES ($1, $2, $3, $4)
        RETURNING id, created_at, updated_at`

	args := []any{wallet.UserID, wallet.Balance, time.Now(), time.Now()}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&wallet.ID, &wallet.CreatedAt, &wallet.UpdatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (m WalletModel) TopUp(walletID int64, amount float64) error {
	query := `
        UPDATE wallets
        SET balance = balance + $1, updated_at = $2
        WHERE id = $3`

	args := []any{amount, time.Now(), walletID}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}

func (m WalletModel) Withdraw(walletID int64, amount float64) error {

	var balance float64
	query := `SELECT balance FROM wallets WHERE id = $1`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, walletID).Scan(&balance)
	if err != nil {
		return err
	}

	if balance < amount {
		return errors.New("insufficient balance")
	}

	query = `
        UPDATE wallets
        SET balance = balance - $1, updated_at = $2
        WHERE id = $3`

	args := []any{amount, time.Now(), walletID}

	_, err = m.DB.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}

func (m WalletModel) GetByID(walletID int64) (*Wallet, error) {
	query := `SELECT id, user_id, balance, created_at, updated_at FROM wallets WHERE id = $1`

	var wallet Wallet

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, walletID).Scan(
		&wallet.ID,
		&wallet.UserID,
		&wallet.Balance,
		&wallet.CreatedAt,
		&wallet.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("wallet not found")
		}
		return nil, err
	}

	return &wallet, nil
}

func (m WalletModel) GetByUserID(userID int64) (*Wallet, error) {
	query := `SELECT id, user_id, balance, created_at, updated_at FROM wallets WHERE user_id = $1`

	var wallet Wallet

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, userID).Scan(
		&wallet.ID,
		&wallet.UserID,
		&wallet.Balance,
		&wallet.CreatedAt,
		&wallet.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("wallet not found")
		}
		return nil, err
	}

	return &wallet, nil
}
