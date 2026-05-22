package postgresql

import (
	"context"
	"database/sql"
	"log/slog"
	"time"

	"github.com/IvanDrf/analyse-search-requests/internal/domain/models"
)

type messageRepo struct {
	db *sql.DB
}

func NewMessageRepo(db *sql.DB) *messageRepo {
	return &messageRepo{
		db: db,
	}
}

func (r *messageRepo) Close() {
	r.db.Close()
}

func (r *messageRepo) SaveMessage(ctx context.Context, message *models.Message) error {
	const query = `
		INSERT INTO searches(search, date, status) VALUES ($1, $2, true) 
		ON CONFLICT (search) DO UPDATE SET search_amount = search_amount + 1
	`

	res, err := r.db.ExecContext(ctx, query, message.SearchMessage, message.Date)
	if _, e := res.RowsAffected(); e != nil {
		slog.Error("any rows weren't affected by inserting new search", slog.String("error", e.Error()))
		return models.Error{
			Message: "no rows were affected by inserting new search",
			Code:    models.ErrCodeInternal,
		}
	}

	if err != nil {
		slog.Error("can't add new search", slog.String("error", err.Error()))
		return models.Error{
			Message: "can't insert new search",
			Code:    models.ErrCodeInternal,
		}
	}

	return nil
}

func (r *messageRepo) FindMostPopularSearches(ctx context.Context, limit uint16, currentTime time.Time) ([]*models.Message, error) {
	const query = `
		SELECT search, search_amount FROM searches 
		WHERE status = true AND EXTRACT(MINUTE FROM $1 - date) = 5
		ORDER BY search_amount DESC
		LIMIT $2
	`

	rows, err := r.db.QueryContext(ctx, query, currentTime, limit)
	if err != nil {
		slog.Error("can't find most popular searches", slog.String("error", err.Error()))
		return nil, models.Error{
			Message: "can't find most popular searches",
			Code:    models.ErrCodeInternal,
		}
	}

	searches := make([]*models.Message, 0, limit)
	for rows.Next() {
		message := &models.Message{}

		if err := rows.Scan(&message.SearchMessage, &message.SearchAmount); err != nil {
			slog.Error("can't scan message from database", slog.String("error", err.Error()))
			return nil, models.Error{
				Message: "can't scan message from database",
				Code:    models.ErrCodeInternal,
			}
		}

		searches = append(searches, message)
	}

	return searches, nil
}

func (r *messageRepo) UpdateSearchStatus(ctx context.Context, status models.MessageStatus, badWord *models.BadWord) error {
	const query = `
		UPDATE searches SET status = $1
		WHERE search LIKE %$2%
	`

	_, err := r.db.ExecContext(ctx, query, status, badWord)
	if err != nil {
		slog.Error("can't update status for searches", slog.String("badWord", badWord.Word))
		return models.Error{
			Message: "can't update status for searches",
			Code:    models.ErrCodeInternal,
		}
	}

	return nil
}
