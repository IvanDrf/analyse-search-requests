package service

import (
	"context"

	"github.com/IvanDrf/analyse-search-requests/internal/domain/models"
)

type BadWordService interface {
	SaveBadWords(ctx context.Context, words []*models.BadWord) error
}
