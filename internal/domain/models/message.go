package models

import (
	"time"

	"github.com/google/uuid"
)

type MessageStatus bool

const (
	ALLOWED   MessageStatus = true
	FORBIDDEN MessageStatus = false
)

type Message struct {
	ID            uuid.UUID `json:"message_id"`
	SearchMessage string    `json:"search_message"`
	Date          time.Time `json:"search_time"`
}

type SearchMessage struct {
	SearchMessage string `json:"search_message"`
	Amount        uint16 `json:"amount"`
}
