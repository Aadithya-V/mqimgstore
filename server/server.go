package server

import (
	"context"
	"database/sql"

	"github.com/Aadithya-V/mqimgstore/queue"
	"github.com/gin-gonic/gin"
)

var (
	ListenAddr = "localhost:8080"
)

type Server struct {
	ctx    context.Context
	router *gin.Engine
	db     *sql.DB
	msgB   *queue.MessageBroker
}

// Function initialises a server,
// maps the routes and returns a pointer to *gin.Engine
func NewServer(db *sql.DB, msgB *queue.MessageBroker) *gin.Engine {
	ctx := context.Background()
	router := gin.Default()
	s := Server{
		ctx,
		router,
		db,
		msgB,
	}

	// Set CORS Policy ?
	// Set health, etc

	s.Product()

	return router
}
