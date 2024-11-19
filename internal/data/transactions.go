package data

import (
	"context"
	"database/sql"
	"time"
)

type Transaction struct {
	ID              int64     `json:"id"`
	WalletID        int64     `json:"wallet_id"`
	Amount          float64   `json:"amount"`
	TransactionType string    `json:"transaction_type"` // "topup" or "withdraw"
	CreatedAt       time.Time `json:"created_at"`
}

type TransactionModel struct {
	DB *sql.DB
}

// RecordTransaction records a new transaction in the transactions table.
func (m TransactionModel) RecordTransaction(walletID int64, amount float64, transactionType string) error {
	query := `
        INSERT INTO transactions (wallet_id, amount, transaction_type, created_at)
        VALUES ($1, $2, $3, $4)`

	args := []any{walletID, amount, transactionType, time.Now()}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}

func (m TransactionModel) GetTransactionsByWalletID(walletID int64) ([]*Transaction, error) {
	query := `SELECT id, wallet_id, amount, transaction_type, created_at FROM transactions WHERE wallet_id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, walletID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []*Transaction

	for rows.Next() {
		var transaction Transaction
		err := rows.Scan(
			&transaction.ID,
			&transaction.WalletID,
			&transaction.Amount,
			&transaction.TransactionType,
			&transaction.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, &transaction)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return transactions, nil
}
