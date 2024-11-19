package main

import (
	"errors"
	"math/rand/v2"
	"net/http"
	"strings"

	"paper.wallet.net/internal/data"
)

func (app *application) topupWallet(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Amount float64 `json:"amount"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		app.invalidAuthenticationTokenResponse(w, r)
		return
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")

	userID, err := app.models.Tokens.Verify(token)
	if err != nil {
		app.invalidAuthenticationTokenResponse(w, r)
		return
	}

	if input.Amount <= 0 {
		app.failedValidationResponse(w, r, map[string]string{"amount": "must be a positive number"})
		return
	}

	user, err := app.models.Users.GetByID(userID)
	if err != nil {
		if errors.Is(err, data.ErrRecordNotFound) {
			app.notFoundResponse(w, r)
		} else {
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	wallet, err := app.models.Wallets.GetByUserID(user.ID)
	if err != nil {
		if errors.Is(err, data.ErrRecordNotFound) {
			app.notFoundResponse(w, r)
		} else {
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// Simple mock payment gateway
	if !mockPaymentGateway(input.Amount) {
		app.serverErrorResponse(w, r, errors.New("payment failed"))
		return
	}

	err = app.models.Wallets.TopUp(wallet.ID, input.Amount)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.models.Transactions.RecordTransaction(wallet.ID, input.Amount, "topup")
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Wallet topped up successfully"))
}

func (app *application) withdrawWallet(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Amount float64 `json:"amount"`
		Pin    string  `json:"pin"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		app.invalidAuthenticationTokenResponse(w, r)
		return
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")

	userID, err := app.models.Tokens.Verify(token)
	if err != nil {
		app.invalidAuthenticationTokenResponse(w, r)
		return
	}

	if input.Amount <= 0 {
		app.failedValidationResponse(w, r, map[string]string{"amount": "must be a positive number"})
		return
	}

	user, err := app.models.Users.GetByID(userID)
	if err != nil {
		if errors.Is(err, data.ErrRecordNotFound) {
			app.notFoundResponse(w, r)
		} else {
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	if user.Pin != input.Pin {
		app.invalidCredentialsResponse(w, r)
		return
	}

	wallet, err := app.models.Wallets.GetByUserID(user.ID)
	if err != nil {
		if errors.Is(err, data.ErrRecordNotFound) {
			app.notFoundResponse(w, r)
		} else {
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	if wallet.Balance < input.Amount {
		app.failedValidationResponse(w, r, map[string]string{"amount": "insufficient balance"})
		return
	}

	err = app.models.Wallets.Withdraw(wallet.ID, input.Amount)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	err = app.models.Transactions.RecordTransaction(wallet.ID, input.Amount, "withdraw")
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Withdrawal successful"))
}

func (app *application) getWalletBalance(w http.ResponseWriter, r *http.Request) {

	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		app.invalidAuthenticationTokenResponse(w, r)
		return
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")

	userID, err := app.models.Tokens.Verify(token)
	if err != nil {
		app.invalidAuthenticationTokenResponse(w, r)
		return
	}

	wallet, err := app.models.Wallets.GetByUserID(userID)
	if err != nil {
		if errors.Is(err, data.ErrRecordNotFound) {
			app.notFoundResponse(w, r)
		} else {
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"balance": wallet.Balance}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) getWalletTransactions(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Pin string `json:"pin"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		app.invalidAuthenticationTokenResponse(w, r)
		return
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")

	userID, err := app.models.Tokens.Verify(token)
	if err != nil {
		app.invalidAuthenticationTokenResponse(w, r)
		return
	}

	user, err := app.models.Users.GetByID(userID)
	if err != nil {
		if errors.Is(err, data.ErrRecordNotFound) {
			app.notFoundResponse(w, r)
		} else {
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	if user.Pin != input.Pin {
		app.invalidCredentialsResponse(w, r)
		return
	}

	wallet, err := app.models.Wallets.GetByUserID(user.ID)
	if err != nil {
		if errors.Is(err, data.ErrRecordNotFound) {
			app.notFoundResponse(w, r)
		} else {
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	transactions, err := app.models.Transactions.GetTransactionsByWalletID(wallet.ID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"transactions": transactions}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// mockPaymentGateway simulates a payment gateway.
func mockPaymentGateway(amount float64) bool {
	// Simulate a 90% success rate
	return rand.Float64() < 0.9
}
