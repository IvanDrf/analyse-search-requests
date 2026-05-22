package models

import (
	"time"
)

type MessageStatus bool

const (
	ALLOWED   MessageStatus = true
	FORBIDDEN MessageStatus = false
)

type Message struct {
	SearchMessage string    `json:"search_message"`
	SearchAmount  uint64    `json:"search_amount"`
	Date          time.Time `json:"search_time"`
}
