package models

import (
	"strings"
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
	Amount        uint64 `json:"amount"`
}

func (m *Message) Normalize() {
	m.SearchMessage = strings.ToLower(m.SearchMessage)
}
