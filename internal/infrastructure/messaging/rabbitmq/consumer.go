package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"sync"

	"github.com/IvanDrf/analyse-search-requests/internal/domain/models"
	"github.com/IvanDrf/analyse-search-requests/internal/domain/ports/service"
	"github.com/rabbitmq/amqp091-go"
)

type searchConsumer struct {
	conn  *amqp091.Connection
	ch    *amqp091.Channel
	queue *amqp091.Queue

	searchService service.SearchService

	workers int
}

func NewSearchConsumer(
	conn *amqp091.Connection, ch *amqp091.Channel,
	queue *amqp091.Queue, searchService service.SearchService,
	workers int,
) *searchConsumer {
	return &searchConsumer{
		conn:          conn,
		ch:            ch,
		queue:         queue,
		searchService: searchService,
		workers:       workers,
	}
}

func (c *searchConsumer) Close() {
	c.ch.Close()
	c.conn.Close()

	slog.Info("SearchConsumer:Close", slog.String("status", "successfully closed search consumer"))
}

func (c *searchConsumer) StartReadingMessages(ctx context.Context) error {
	delivery, err := c.ch.Consume(c.queue.Name, "", false, false, false, false, nil)
	if err != nil {
		return models.Error{
			Message: fmt.Sprintf("can't start reading messages, error=%s", err),
		}
	}

	wg := new(sync.WaitGroup)
	for range c.workers {
		wg.Add(1)
		go c.worker(ctx, delivery)
	}

	<-ctx.Done()
	wg.Wait()
	return nil
}

func (c *searchConsumer) worker(ctx context.Context, delivery <-chan amqp091.Delivery) {
	for {
		select {
		case <-ctx.Done():
			return

		case message, ok := <-delivery:
			if !ok {
				return
			}

			if err := c.processMessage(ctx, message); err != nil {
				slog.Error("", slog.String("error", err.Error()))
			}
		}
	}
}

func (c *searchConsumer) processMessage(ctx context.Context, message amqp091.Delivery) error {
	search := models.Message{}
	if err := json.Unmarshal(message.Body, &search); err != nil {
		message.Ack(false)
		return models.Error{
			Message: fmt.Sprintf("can't parse incoming message from queue, error=%s", err),
		}
	}

	if err := c.searchService.SaveSearch(ctx, &search); err != nil {
		message.Nack(false, true)
		return models.Error{
			Message: fmt.Sprintf("can't save new search error=%s", err),
		}
	}

	message.Ack(false)
	return nil
}
