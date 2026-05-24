package service

import (
	"context"

	"github.com/IvanDrf/analyse-search-requests/internal/domain/models"
)

type SearchService interface {
	SaveSearch(ctx context.Context, message *models.Message) error
	FindMostPopularSearches(ctx context.Context, limit int) ([]*models.SearchMessage, error)

	SaveBadWord(ctx context.Context, badWord string) error
	DeleteBadWord(ctx context.Context, badWord string) error

	Close()
}
