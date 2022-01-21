package api

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	db "github.com/harrychopra/go-api/db/models"
)

// Server serves HTTP requests for banking service
type Server struct {
	store  db.Store
	router *gin.Engine
}

// NewServer creates a new HTTP server and sets up routing
func NewServer(store db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()
	// Register custom validator
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}
	router.POST("/accounts", server.CreateAccount)
	router.GET("/accounts/:id", server.GetAccount)
	router.GET("/accounts", server.ListAccounts)
	router.POST("/transfers", server.CreateTransfer)
	server.router = router
	return server
}

// Start starts an HTTP Server and listens on the input "address:port"
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errResponse(err error) *gin.H {
	return &gin.H{"error": err.Error()}
}
