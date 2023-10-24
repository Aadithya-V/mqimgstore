package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/go-sql-driver/mysql"
	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/Aadithya-V/mqimgstore/imgcompressionservice/handlers"
)

// Constants specifying the listening addresses of
var (
	MySQLAddr = "127.0.0.1:3306"
)

func main() {

	db := initSqlDB()

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	// Should be CONSISTENT WITH PRODUCER!!
	q, err := ch.QueueDeclare(
		"img-compression-queue", // name
		true,                    // durable
		false,                   // delete when unused
		false,                   // exclusive
		false,                   // no-wait
		nil,                     // arguments
	)
	failOnError(err, "Failed to declare a queue")

	err = ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	failOnError(err, "Failed to set QoS")

	msgs, err := ch.Consume( // read consume's doc for ack mechanisms
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	var forever chan struct{}

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
			go handlers.ServiceWorker(d, db) // TODO: load balance. Stop consuming if CPU util exceeds certain amt- use benchmarking here?
		}
	}()

	// save compressed img locally (and expose service endpoints?)

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
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

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}
