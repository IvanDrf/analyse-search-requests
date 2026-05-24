package rules

import (
	"time"

	"github.com/IvanDrf/analyse-search-requests/internal/domain/models"
)

type MessageValidator struct {
}

func NewMessageValidator() *MessageValidator {
	return &MessageValidator{}
}

func (v *MessageValidator) ValidateMessage(message *models.Message) error {
	if message.SearchMessage == "" {
		return models.Error{
			Message: "search message is empty",
			Code:    models.ErrCodeInvalidArgument,
		}
	}

	now := time.Now().UTC()
	if !now.After(message.Date) {
		return models.Error{
			Message: "message date is after current time",
			Code:    models.ErrCodeInvalidArgument,
		}
	}

	return nil
}
