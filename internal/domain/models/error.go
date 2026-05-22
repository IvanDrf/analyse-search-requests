package models

type ErrorCode int

const (
	ErrCodeInternal ErrorCode = iota
	ErrCodeInvalidArgument
)

type Error struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}
