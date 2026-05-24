package models

import "fmt"

type ErrorCode int

const (
	ErrCodeInternal ErrorCode = iota
	ErrCodeInvalidArgument
	ErrCodeDuplicateMessage
)

type Error struct {
	Message string    `json:"message"`
	Code    ErrorCode `json:"-"`
}

func (e Error) Error() string {
	return fmt.Sprintf("message: %s, code: %d", e.Message, e.Code)
}
