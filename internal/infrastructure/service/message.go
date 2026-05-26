package service

import (
	"context"
	"log/slog"
	"time"

	"github.com/IvanDrf/analyse-search-requests/internal/domain/models"
	"github.com/IvanDrf/analyse-search-requests/internal/domain/ports/repo"
	"github.com/IvanDrf/analyse-search-requests/internal/domain/rules"
)

type searchService struct {
	timeInterval int

	messageRepo repo.MessageRepo

	validator *rules.MessageValidator
}

func NewSearchService(timeInterval int, repo repo.MessageRepo, validator *rules.MessageValidator) *searchService {
	return &searchService{
		timeInterval: timeInterval,
		messageRepo:  repo,
		validator:    validator,
	}
}

func (s *searchService) Close() {
	s.messageRepo.Close()
}

func (s *searchService) SaveSearch(ctx context.Context, message *models.Message) error {
	slog.Info("SearchService:SaveSearch", slog.String("search", message.SearchMessage))
	if err := s.validator.ValidateMessage(message); err != nil {
		slog.Info("SearchService:SaveSearch", slog.String("error", err.Error()))
		return err
	}

	message.Normalize()

	if err := s.messageRepo.SaveSearch(ctx, message); err != nil {
		slog.Error("SearchService:SaveSearch", slog.String("error", err.Error()))
		return models.Error{
			Message: "can't save message",
			Code:    models.ErrCodeInternal,
		}
	}

	slog.Info("SearchService:SaveSearch", slog.String("search", message.SearchMessage), slog.String("status", "successfully saved search"))
	return nil
}

func (s *searchService) FindMostPopularSearches(ctx context.Context, limit int) ([]*models.SearchMessage, error) {
	slog.Info("SearchService:FindMostPopularSearches", slog.Int("limit", limit))
	if err := rules.ValidateLimit(limit); err != nil {
		slog.Info("SearchService:FindMostPopularSearches", slog.String("error", err.Error()))
		return nil, err
	}

	now := time.Now().UTC()
	searches, err := s.messageRepo.FindMostPopularSearches(ctx, limit, now, s.timeInterval)
	if err != nil {
		slog.Error("SearchService:FindMostPopularSearches", slog.String("error", err.Error()))
		return nil, models.Error{
			Message: "can't find most popular searches",
			Code:    models.ErrCodeInternal,
		}
	}

	slog.Info("SearchService:FindMostPopularSearches", slog.Int("limit", limit), slog.String("status", "successfully found most popular searches"))
	return searches, nil
}

func (s *searchService) SaveBadWord(ctx context.Context, badWord string) error {
	slog.Info("SearchService:SaveBadWord", slog.String("bad_word", badWord))
	if err := rules.ValidateBadWord(badWord); err != nil {
		slog.Info("SearchService:SaveBadWord", slog.String("error", err.Error()))
		return err
	}

	if err := s.messageRepo.SaveBadWord(ctx, badWord); err != nil {
		slog.Error("SearchService:SaveBadWord", slog.String("error", err.Error()))
		return models.Error{
			Message: "can't save bad word",
			Code:    models.ErrCodeInternal,
		}
	}

	slog.Info("SearchService:SaveBadWord", slog.String("bad_word", badWord), slog.String("status", "successfully saved bad word"))
	return nil
}

func (s *searchService) DeleteBadWord(ctx context.Context, badWord string) error {
	slog.Info("SearchService:DeleteBadWord", slog.String("bad_word", badWord))
	if err := rules.ValidateBadWord(badWord); err != nil {
		slog.Info("SearchService:DeleteBadWord", slog.String("error", err.Error()))
		return err
	}

	if err := s.messageRepo.DeleteBadWord(ctx, badWord); err != nil {
		slog.Error("SearchService:DeleteBadWord", slog.String("error", err.Error()))
		return models.Error{
			Message: "can't delete bad word",
			Code:    models.ErrCodeInternal,
		}
	}

	slog.Info("SearchService:DeleteBadWord", slog.String("bad_word", badWord), slog.String("status", "successfully deleted bad word"))
	return nil
}
