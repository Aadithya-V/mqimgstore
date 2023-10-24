package queue

import (
	"context"

	amqp "github.com/rabbitmq/amqp091-go"
)

type MessageBroker struct {
	Connection *amqp.Connection
	Channel    *amqp.Channel
}

const (
	ImageCompressionQueue = "img-compression-queue"
)

func NewMessageBroker() *MessageBroker {
	// var s *MessageBroker // return s=nil, error in the future if api not requiring message broker are added to this go service and dont call log.Panic

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ") // MUST PANIC?  -- restarts pod ;
	// close conn and ch in main using *MessageBroker.Close()

	ch, err := conn.Channel() // TODO: check if the deferred mqConn.close at main also closes mqCh at the end of main to be safer.
	failOnError(err, "Failed to open a channel")

	err = ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)

	err = declareQueues(ch)
	failOnError(err, "Failed to declare a queue")

	return &MessageBroker{Connection: conn, Channel: ch}
}

func declareQueues(ch *amqp.Channel) error {
	_, err := ch.QueueDeclare( // _ == q
		ImageCompressionQueue, // name
		true,                  // durable
		false,                 // delete when unused
		false,                 // exclusive
		false,                 // no-wait
		nil,                   // arguments
	)
	return err // future return []q
}

func (s *MessageBroker) Close() {
	s.Channel.Close()
	s.Connection.Close()
}

func (s *MessageBroker) Publish(ctx context.Context, queueName string, msgBody []byte) error {

	return s.Channel.PublishWithContext(ctx,
		"",        // exchange
		queueName, // routing key
		true,      // mandatory // TODO: add NotifyReturn
		false,     // immediate
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         msgBody,
		})
}

func (s *MessageBroker) Consume(queueName string) (<-chan amqp.Delivery, error) {
	return s.Channel.Consume(
		queueName, // queue
		"",        // consumer
		false,     // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
}
