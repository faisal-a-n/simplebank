package api

import (
	"database/sql"
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

type userDetailsResponse struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	CreatedAt int64  `json:"created_at"`
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

	response := userResponseBuilder(user)
	ctx.JSON(http.StatusCreated, responseHandler(200, "User created", response))
}

type loginUserRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password,onempty" binding:"required,min=6"`
}

type loginReponse struct {
	AccessToken string              `json:"access_token"`
	User        userDetailsResponse `json:"user"`
}

func (server *Server) loginUser(ctx *gin.Context) {
	var req loginUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	user, err := server.store.GetUserByEmail(ctx, req.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(errors.New("Email is not registered")))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	if err := util.CheckPassword(req.Password, user.Password); err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(errors.New("Invalid password entered")))
		return
	}
	access_token, err := server.tokenMaker.CreateToken(user.ID, server.config.ACCESS_TOKEN_DURATION)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	response := loginReponse{
		AccessToken: access_token,
		User:        userResponseBuilder(user),
	}
	ctx.JSON(http.StatusOK, responseHandler(200, "You have logged in successfully", response))
}

func userResponseBuilder(user db.User) userDetailsResponse {
	return userDetailsResponse{
		Name:      user.Name,
		Email:     user.Email,
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
	}
}
