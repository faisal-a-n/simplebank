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
	"github.com/lib/pq"
)

type createAccountRequest struct {
	Name     string `json:"name" binding:"required"`
	Currency string `json:"currency" binding:"required,currency"`
}

type getAccountReq struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

type getAccountsReq struct {
	PageID int32 `form:"page_id" binding:"required,min=1"`
	Count  int32 `form:"count" binding:"required,min=5,max=15"`
}

func (server *Server) createAccount(ctx *gin.Context) {
	var req createAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authPayloadKey).(*token.Payload)
	arg := db.CreateAccountParams{
		Name:      req.Name,
		UserID:    authPayload.UserID,
		Currency:  req.Currency,
		Balance:   0,
		CreatedAt: time.Now().Unix(),
	}

	account, err := server.store.CreateAccount(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				err = fmt.Errorf("User [%d] already has [%s] currency account", authPayload.UserID, req.Currency)
			case "foreign_key_violation":
				err = fmt.Errorf("User [%d] doesn't exist", authPayload.UserID)
			}
			ctx.JSON(http.StatusForbidden, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusCreated, responseHandler(200, "Account has been created", account))
}

//Get account by id
func (server *Server) getAccount(ctx *gin.Context) {
	var req getAccountReq
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("Invalid ID")))
		return
	}

	authPayload := ctx.MustGet(authPayloadKey).(*token.Payload)

	account, err := server.store.GetAccount(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(errors.New("No account with this id")))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	if account.UserID != authPayload.UserID {
		ctx.JSON(http.StatusUnauthorized, errorResponse(errors.New("Account does not belong to the user")))
		return
	}
	ctx.JSON(http.StatusOK, responseHandler(200, "Data fetched", account))
}

//Get accounts list
func (server *Server) getAccounts(ctx *gin.Context) {
	var req getAccountsReq
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authPayloadKey).(*token.Payload)

	accounts, err := server.store.ListAccountsForUser(ctx, db.ListAccountsForUserParams{
		UserID: authPayload.UserID,
		Limit:  req.Count,
		Offset: (req.PageID - 1) * req.Count,
	})

	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(errors.New("No accounts available")))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, responseHandler(200, "Data fetched", accounts))
}
