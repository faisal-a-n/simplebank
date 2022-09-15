package api

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"time"

	db "github.com/faisal-a-n/simplebank/db/sqlc"
	"github.com/faisal-a-n/simplebank/token"
	"github.com/gin-gonic/gin"
)

type createTransferRequest struct {
	FromAccountID int64  `json:"from_account_id" binding:"required,min=1"`
	ToAccountID   int64  `json:"to_account_id" binding:"required,min=1"`
	Amount        int64  `json:"amount" binding:"required,min=1"`
	Currency      string `json:"currency" binding:"required,currency"`
}

type transferEntries struct {
	fromEntry db.Entry
	toEntry   db.Entry
}

func (server *Server) createTransfer(ctx *gin.Context) {
	var req createTransferRequest
	var err error
	if err = ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	fromAccount, fromCheck := server.checkCurrency(ctx, req.FromAccountID, req.Currency)
	_, toCheck := server.checkCurrency(ctx, req.ToAccountID, req.Currency)

	if !(fromCheck && toCheck) {
		return
	}

	if !checkOwnershipAndBalance(ctx, fromAccount, req.Amount) {
		return
	}

	fromEntryParams := buildEntryParams(req.FromAccountID, -req.Amount)
	toEntryParams := buildEntryParams(req.ToAccountID, req.Amount)
	entries := transferEntries{}

	entries.fromEntry, err = server.store.CreateEntry(ctx, fromEntryParams)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	entries.toEntry, err = server.store.CreateEntry(ctx, toEntryParams)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.TransferTxParams{
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		FromEntryID:   entries.fromEntry.ID,
		ToEntryID:     entries.toEntry.ID,
		Amount:        req.Amount,
		CreatedAt:     time.Now().Unix(),
	}

	transaction, err := server.store.TransferTx(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusCreated, responseHandler(200, "Transaction has been made", transaction))
}

func buildEntryParams(accountID int64, amount int64) db.CreateEntryParams {
	return db.CreateEntryParams{
		AccountID: accountID,
		Amount:    amount,
		CreatedAt: time.Now().Unix(),
	}
}

func (server *Server) checkCurrency(ctx *gin.Context, accountID int64, currency string) (db.Account, bool) {
	account, err := server.store.GetAccount(ctx, accountID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(fmt.Errorf("Provided account [%d] doesn't exist", accountID)))
			return db.Account{}, false
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return db.Account{}, false
	}
	if account.Currency != currency {
		err := fmt.Errorf("account [%d] currency mismatch: [%s] vs [%s]", accountID, account.Currency, currency)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return db.Account{}, false
	}

	return account, true
}

func checkOwnershipAndBalance(ctx *gin.Context, account db.Account, amount int64) bool {
	authPayload := ctx.MustGet(authPayloadKey).(*token.Payload)

	if account.UserID != authPayload.UserID {
		ctx.JSON(http.StatusUnauthorized, errorResponse(errors.New("Account does not belong to the user")))
		return false
	}

	if account.Balance < amount {
		ctx.JSON(http.StatusUnauthorized, errorResponse(errors.New("Account does not have enough balance")))
		return false
	}
	return true
}
