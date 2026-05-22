package repo

import (
	"context"

	"github.com/IvanDrf/analyse-search-requests/internal/domain/models"
)

type BadWordRepo interface {
	SaveBadWords(ctx context.Context, words []*models.BadWord) error

	Close()
}
