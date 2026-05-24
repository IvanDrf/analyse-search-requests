package rabbitmq

import (
	"log"

	"github.com/IvanDrf/analyse-search-requests/internal/config"
	"github.com/rabbitmq/amqp091-go"
)

func Connect(config *config.RabbitMQConfig) *amqp091.Connection {
	conn, err := amqp091.Dial(config.DSN())
	if err != nil {
		log.Fatalf("can't connect to rabbitmq, error=%s", err)
	}

	return conn
}
