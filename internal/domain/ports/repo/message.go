package repo

import (
	"context"

	"github.com/IvanDrf/analyse-search-requests/internal/domain/models"
)

type MessageRepo interface {
	SaveMessage(ctx context.Context, message *models.Message) error
	FindMostPopularSearches(ctx context.Context, limit uint16) ([]*models.Message, error)
}
