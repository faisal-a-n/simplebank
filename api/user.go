package api

import (
	"errors"
	"net/http"
	"time"

	db "github.com/faisal-a-n/simplebank/db/sqlc"
	"github.com/faisal-a-n/simplebank/util"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type createUserRequest struct {
	Name     string `json:"name" binding:"required,alpha,min=6"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password,onempty" binding:"required,min=6"`
}

func (server *Server) createUser(ctx *gin.Context) {
	var req createUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	hash, err := util.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	args := db.CreateUserParams{
		Name:              req.Name,
		Email:             req.Email,
		Password:          hash,
		PasswordChangedAt: time.Now().Unix(),
		CreatedAt:         time.Now().Unix(),
	}
	user, err := server.store.CreateUser(ctx, args)
	if err != nil {
		if pqError, ok := err.(*pq.Error); ok {
			if pqError.Code.Name() == "unique_violation" {
				ctx.JSON(http.StatusInternalServerError, errorResponse(errors.New("User with the same email already exists")))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	response := struct {
		ID        int64  `json:"id"`
		Name      string `json:"name"`
		Email     string `json:"email"`
		CreatedAt int64  `json:"created_at"`
	}{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}
	ctx.JSON(http.StatusCreated, responseHandler(200, "User created", response))
}
