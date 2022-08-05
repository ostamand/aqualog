package api

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/ostamand/aqualog/storage"
	"github.com/ostamand/aqualog/token"
	"github.com/ostamand/aqualog/util"
)

const ginReleaseEnv = "GIN_MODE"

type Server struct {
	config     util.Config
	storage    storage.Storage
	tokenMaker token.TokenMaker
	router     *gin.Engine
}

func NewServer(config util.Config, s storage.Storage) *Server {
	// set production mode if needed
	mode := os.Getenv(ginReleaseEnv)
	if mode != "" {
		gin.SetMode(gin.ReleaseMode)
	}

	t, err := token.NewPasetoMaker(config.TokenKey)
	if err != nil {
		log.Fatalf("could not initialize authentication: %s", errorResponse(err))
	}

	server := &Server{storage: s, tokenMaker: t, config: config}

	router := gin.Default()

	router.POST("/users", server.createUser)
	router.POST("/login", server.login)

	authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))
	authRoutes.POST("/params", server.createParam)
	authRoutes.GET("/params", server.getParams)

	server.router = router
	return server
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
