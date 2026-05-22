package models

import (
	"time"
)

type Message struct {
	SearchMessage string    `json:"search_message"`
	Date          time.Time `json:"search_time"`
}
