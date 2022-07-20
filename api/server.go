package api

import (
	"github.com/gin-gonic/gin"
	"github.com/ostamand/aqualog/storage"
)

type Server struct {
	storage storage.Storage
	router  *gin.Engine
}

func NewServer(s storage.Storage) *Server {
	server := &Server{storage: s}
	router := gin.Default()

	router.POST("/values", server.saveValue)
	router.GET("/values/:id", server.getValue)

	server.router = router
	return server
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
