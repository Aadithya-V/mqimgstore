package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/Aadithya-V/mqimgstore/database"
	"github.com/Aadithya-V/mqimgstore/queue"
	"github.com/Aadithya-V/mqimgstore/server"
	"github.com/gin-gonic/gin"
)

// Constants specifying the listening addresses of
// the MySQL server and the gin router engine.
var (
	ListenAddr = "localhost:8080"
)

/* type closable interface {
	Close()
}

// connections made by functions that need to be closed only at the end of main()
var closeCh = make(chan closable, 64) // Large value to prevent enqueue deadlock */

func main() {

	// Initialize database connection
	db := database.NewMySQLSession()

	// Connect to RabbitMQ
	msgB := queue.NewMessageBroker()
	defer msgB.Close()

	// Initialize router gin engine
	router := server.NewServer(db, msgB)

	startHttpServer(router) // blocking. Spin this off into a go routine if subsequent code is added.

}

func startHttpServer(router *gin.Engine) {
	srv := &http.Server{
		Addr:    ListenAddr,
		Handler: router,
	}

	idleConnsClosed := make(chan struct{})

	go func() {
		// for graceful shutdown..
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		// We received an interrupt signal, shut down.
		log.Println("Shutting down server...")
		if err := srv.Shutdown(context.Background()); err != nil {
			// Error from closing listener(s), or context timeout
			log.Printf("HTTP server Shutdown: %v", err)
		}
		close(idleConnsClosed)
	}()

	// blocking service of connections
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("HTTP server ListenAndServe: %v", err)
	}

	<-idleConnsClosed
}
