package rabbitmq

import (
	"log"

	"github.com/IvanDrf/analyse-search-requests/internal/config"
	"github.com/rabbitmq/amqp091-go"
)

func Connect(config *config.RabbitMQConfig) (*amqp091.Connection, *amqp091.Channel, *amqp091.Queue) {
	conn, err := amqp091.Dial(config.DSN())
	if err != nil {
		log.Fatalf("can't connect to rabbitmq, error=%s", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		log.Fatalf("can't declare channel, error=%s", err)
	}

	queue, err := ch.QueueDeclare(config.Queue, true, false, false, false, nil)
	if err != nil {
		ch.Close()
		conn.Close()
		log.Fatalf("can't declare queue, error=%s", err)
	}

	return conn, ch, &queue
}
