package api

import (
	db "github.com/faisal-a-n/simplebank/db/sqlc"
	"github.com/gin-gonic/gin"
)

//Serves http requests for banking service
type Server struct {
	store  db.Store
	router *gin.Engine
}

//Create new server and setup routing
func NewServer(store db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	//add routes to router
	router.POST("/accounts", server.createAccount)
	router.GET("/accounts/:id", server.getAccount)
	router.GET("/accounts", server.getAccounts)

	server.router = router
	return server
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{
		"message": "There was an error",
		"error":   err.Error(),
	}
}

func responseHandler(code int, message string, data interface{}) gin.H {
	return gin.H{
		"code":    code,
		"message": message,
		"data":    data,
	}
}
