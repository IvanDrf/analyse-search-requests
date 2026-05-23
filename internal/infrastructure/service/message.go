package service

import (
	"context"

	"github.com/IvanDrf/analyse-search-requests/internal/domain/models"
	"github.com/IvanDrf/analyse-search-requests/internal/domain/ports/repo"
	"github.com/IvanDrf/analyse-search-requests/internal/domain/rules"
)

type messageService struct {
	messageRepo repo.MessageRepo

	validator *rules.MessageValidator
}

func NewMessageService(repo repo.MessageRepo, validator *rules.MessageValidator) *messageService {
	return &messageService{
		messageRepo: repo,
		validator:   validator,
	}
}

func (s *messageService) SaveMessage(ctx context.Context, message *models.Message) error {
	if err := s.validator.ValidateMessage(message); err != nil {
		return err
	}

	if err := s.messageRepo.SaveMessage(ctx, message); err != nil {
		return err
	}

	return nil
}
