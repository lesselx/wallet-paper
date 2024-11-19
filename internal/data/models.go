package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
)

type Models struct {
	Users        UserModel
	Tokens       TokenModel
	Wallets      WalletModel
	Transactions TransactionModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Users:        UserModel{DB: db},
		Tokens:       TokenModel{DB: db},
		Wallets:      WalletModel{DB: db},
		Transactions: TransactionModel{DB: db},
	}
}
