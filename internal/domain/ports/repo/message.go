package repo

import (
	"context"
	"time"

	"github.com/IvanDrf/analyse-search-requests/internal/domain/models"
)

type MessageRepo interface {
	SaveSearch(ctx context.Context, message *models.Message) error
	SaveBadWord(ctx context.Context, badWord string) error
	DeleteBadWord(ctx context.Context, badWord string) error

	FindMostPopularSearches(ctx context.Context, limit int, start time.Time, interval int) ([]*models.SearchMessage, error)

	Close()
}
