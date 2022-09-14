package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	db "github.com/faisal-a-n/simplebank/db/sqlc"
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

	if !(server.checkCurrency(ctx, req.FromAccountID, req.Currency) &&
		server.checkCurrency(ctx, req.ToAccountID, req.Currency)) {
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

func (server *Server) checkCurrency(ctx *gin.Context, accountID int64, currency string) bool {
	account, err := server.store.GetAccount(ctx, accountID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return false
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return false
	}
	if account.Currency != currency {
		err := fmt.Errorf("account [%d] currency mismatch: [%s] vs [%s]", accountID, account.Currency, currency)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return false
	}
	return true
}
