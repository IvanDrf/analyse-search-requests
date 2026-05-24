package rules

import "github.com/IvanDrf/analyse-search-requests/internal/domain/models"

func ValidateLimit(limit int) error {
	if limit <= 0 {
		return models.Error{
			Message: "limit must be positive",
			Code:    models.ErrCodeInvalidArgument,
		}
	}

	return nil
}
