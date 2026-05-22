package rules

import (
	"time"

	"github.com/IvanDrf/analyse-search-requests/internal/domain/models"
)

type MessageValidator struct {
	location *time.Location
}

func NewMessageValidator(location *time.Location) *MessageValidator {
	return &MessageValidator{
		location: location,
	}
}

func (v *MessageValidator) ValidateMessage(message *models.Message) error {
	if message.SearchMessage == "" {
		return models.Error{
			Message: "search message is empty",
			Code:    models.ErrCodeInvalidArgument,
		}
	}

	now := time.Now().In(v.location)
	if !now.After(message.Date) {
		return models.Error{
			Message: "message date is after current time",
			Code:    models.ErrCodeInvalidArgument,
		}
	}

	return nil
}
