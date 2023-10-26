package main

import (
	"log"

	"github.com/Aadithya-V/mqimgstore/database"
	"github.com/Aadithya-V/mqimgstore/imgcompressionservice/handlers"
	"github.com/Aadithya-V/mqimgstore/queue"
)

// Constants specifying the listening addresses of
var (
	MySQLAddr = "127.0.0.1:3306"
)

func main() {

	db := database.NewMySQLSession()

	msgQ := queue.NewMessageBroker()
	defer msgQ.Close()

	msgs, err := msgQ.Consume(queue.ImageCompressionQueue)

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

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}
