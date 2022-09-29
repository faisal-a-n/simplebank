package api

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type renewTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type renewTokenResponse struct {
	AccessToken       string    `json:"access_token"`
	AccessTokenExpiry time.Time `json:"access_token_expires_at"`
}

func (server *Server) renewToken(ctx *gin.Context) {
	var req renewTokenRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	refreshToken, err := server.tokenMaker.VerifyToken(req.RefreshToken)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	if time.Now().After(refreshToken.ExpiredAt) {
		ctx.JSON(http.StatusUnauthorized, errorResponse(errors.New("Refresh token has expired")))
		return
	}

	session, err := server.store.GetSession(ctx, refreshToken.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(errors.New("Session is invalid")))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	if session.IsBlocked {
		ctx.JSON(http.StatusUnauthorized, errorResponse(fmt.Errorf("Refresh token is invalid")))
		return
	}
	if session.UserID != refreshToken.UserID {
		ctx.JSON(http.StatusUnauthorized, errorResponse(fmt.Errorf("Incorrect session user")))
		return
	}
	if session.RefreshToken != req.RefreshToken {
		ctx.JSON(http.StatusUnauthorized, errorResponse(fmt.Errorf("Refresh token mismatch")))
		return
	}

	access_token, payload, err := server.tokenMaker.CreateToken(refreshToken.UserID, server.config.ACCESS_TOKEN_DURATION)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	response := renewTokenResponse{
		AccessToken:       access_token,
		AccessTokenExpiry: payload.ExpiredAt,
	}
	ctx.JSON(http.StatusCreated, responseHandler(200, "Token refreshed", response))
}
