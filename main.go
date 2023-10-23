package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"

	"database/sql"

	"github.com/Aadithya-V/mqimgstore/handlers"
	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	amqp "github.com/rabbitmq/amqp091-go"
)

// Constants specifying the listening addresses of
// the MySQL server and the gin router engine.
var (
	ListenAddr = "localhost:8080"
	MySQLAddr  = "127.0.0.1:3306"
)

/* type closable interface {
	Close()
}

// connections made by functions that need to be closed only at the end of main()
var closeCh = make(chan closable, 64) // Large value to prevent enqueue deadlock */

func main() {

	// Initialize database connection
	db := initSqlDB()

	// Connect to RabbitMq
	mqConn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ") // MUST PANIC?  -- restarts pod ;
	defer mqConn.Close()

	// Initialize router gin engine
	router := initRouter(db, mqConn)

	startHttpServer(router) // blocking. Spin this off into a go routine if subsequent code is added.

}

// Function initRouter() initialises a router,
// maps the routes and returns a pointer to it
// which is *gin.Engine
func initRouter(db *sql.DB, mqConn *amqp.Connection) *gin.Engine {
	router := gin.Default()

	router.POST("/product", handlers.AddProduct(db, mqConn))

	return router
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

func initSqlDB() *sql.DB {
	// Capture connection properties.
	cfg := mysql.Config{
		User:   "go",  //os.Getenv("DBUSER"),
		Passwd: "123", //os.Getenv("DBPASS"),
		Net:    "tcp",
		Addr:   MySQLAddr,
		DBName: "test",
	}
	// Get a database handle.
	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Println("MySQL Connected!")

	return db
}
