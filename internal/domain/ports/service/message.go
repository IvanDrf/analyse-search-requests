package service

import (
	"context"

	"github.com/IvanDrf/analyse-search-requests/internal/domain/models"
)

type MessageService interface {
	SaveSearch(ctx context.Context, message *models.Message) error
	FindMostPopularSearches(ctx context.Context, amount uint16) ([]*models.Message, error)
}
