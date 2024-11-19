package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {

	router := httprouter.New()

	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)
	router.HandlerFunc(http.MethodPost, "/v1/users", app.registerUserHandler)
	router.HandlerFunc(http.MethodPost, "/v1/tokens/authentication", app.createAuthenticationTokenHandler)

	router.HandlerFunc(http.MethodPost, "/v1/wallets/topup", app.requireAuthenticatedUser(app.topupWallet))
	router.HandlerFunc(http.MethodPost, "/v1/wallets/withdraw", app.requireAuthenticatedUser(app.withdrawWallet))
	router.HandlerFunc(http.MethodGet, "/v1/wallets/balance", app.requireAuthenticatedUser(app.getWalletBalance))
	router.HandlerFunc(http.MethodGet, "/v1/wallets/transactions", app.requireAuthenticatedUser(app.getWalletTransactions))

	return app.recoverPanic(app.authenticate(router))

}
