package messaging

import "context"

type MessageConsumer interface {
	StartReadingMessages(ctx context.Context) error

	Close()
}
