package api

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/ostamand/aqualog/storage"
	"github.com/ostamand/aqualog/token"
	"github.com/ostamand/aqualog/util"
)

type Server struct {
	config     util.Config
	storage    storage.Storage
	tokenMaker token.TokenMaker
	router     *gin.Engine
}

func NewServer(config util.Config, s storage.Storage) *Server {
	gin.SetMode(config.Mode)

	t, err := token.NewPasetoMaker(config.TokenKey)
	if err != nil {
		log.Fatalf("could not initialize authentication: %s", errorResponse(err))
	}

	server := &Server{storage: s, tokenMaker: t, config: config}

	router := gin.Default()
	router.Use(corsMiddleware())

	router.GET("/data", server.getData)
	router.POST("/users", server.createUser)
	router.POST("/login", server.login)
	router.POST("/renew_token", server.renewToken)

	authRoutes := router.Use(authMiddleware(server.tokenMaker))
	authRoutes.POST("/params", server.createParam)
	authRoutes.GET("/params", server.getParams)
	authRoutes.GET("/params/:id", server.getParam)
	authRoutes.GET("/params/summary", server.getSummary)

	server.router = router

	return server
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
