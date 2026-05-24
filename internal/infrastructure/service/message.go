package service

import (
	"context"
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

func NewMessageService(timeInterval int, repo repo.MessageRepo, validator *rules.MessageValidator) *searchService {
	return &searchService{
		timeInterval: timeInterval,
		messageRepo:  repo,
		validator:    validator,
	}
}

func (s *searchService) SaveSearch(ctx context.Context, message *models.Message) error {
	if err := s.validator.ValidateMessage(message); err != nil {
		return err
	}

	if err := s.messageRepo.SaveSearch(ctx, message); err != nil {
		return models.Error{
			Message: "can't save message",
			Code:    models.ErrCodeInternal,
		}
	}

	return nil
}

func (s *searchService) FindMostPopularSearches(ctx context.Context, limit int) ([]*models.SearchMessage, error) {
	if err := rules.ValidateLimit(limit); err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	searches, err := s.messageRepo.FindMostPopularSearches(ctx, limit, now, s.timeInterval)
	if err != nil {
		return nil, models.Error{
			Message: "can't find most popular searches",
			Code:    models.ErrCodeInternal,
		}
	}

	return searches, nil
}

func (s *searchService) SaveBadWord(ctx context.Context, badWord string) error {
	if err := rules.ValidateBadWord(badWord); err != nil {
		return err
	}

	if err := s.messageRepo.SaveBadWord(ctx, badWord); err != nil {
		return models.Error{
			Message: "can't save bad word",
			Code:    models.ErrCodeInternal,
		}
	}

	return nil
}

func (s *searchService) DeleteBadWord(ctx context.Context, badWord string) error {
	if err := rules.ValidateBadWord(badWord); err != nil {
		return err
	}

	if err := s.messageRepo.DeleteBadWord(ctx, badWord); err != nil {
		return models.Error{
			Message: "can't delete bad word",
			Code:    models.ErrCodeInternal,
		}
	}

	return nil
}
