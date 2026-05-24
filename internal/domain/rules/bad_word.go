package rules

import "github.com/IvanDrf/analyse-search-requests/internal/domain/models"

func ValidateBadWord(badWord string) error {
	if badWord == "" {
		return models.Error{
			Message: "bad word can't be empty",
			Code:    models.ErrCodeInvalidArgument,
		}
	}

	return nil
}
