package postgresql

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/IvanDrf/analyse-search-requests/internal/domain/models"
	"github.com/lib/pq"
)

type badWordRepo struct {
	db *sql.DB
}

func NewBadWordRepo(db *sql.DB) *badWordRepo {
	return &badWordRepo{
		db: db,
	}
}

func (r *badWordRepo) Close() {
	r.db.Close()
}

func (r *badWordRepo) SaveBadWords(ctx context.Context, words []*models.BadWord) error {
	const query = `
		INSERT INTO bad_words(word) VALUES($1) ON CONFLICT DO NOTHING
	`
	args := createArgsFromWords(words)

	_, err := r.db.ExecContext(ctx, query, pq.Array(args))

	if err != nil {
		slog.Error("can't save new bad words")
		return models.Error{
			Message: "can't save new bad words",
			Code:    models.ErrCodeInternal,
		}
	}

	return nil
}

func (r *badWordRepo) DeleteBadWords(ctx context.Context, words []*models.BadWord) error {
	const query = `
		DELETE FROM bad_words WHERE word IN ($1)
	`

	args := createArgsFromWords(words)

	_, err := r.db.ExecContext(ctx, query, args)
	if err != nil {
		slog.Error("can't delete bad words", slog.String("error", err.Error()))
		return models.Error{
			Message: "can't delete bad words",
			Code:    models.ErrCodeInternal,
		}
	}

	return nil
}

func createArgsFromWords(words []*models.BadWord) []string {
	args := make([]string, 0, len(words))

	for i := range words {
		args = append(args, words[i].Word)
	}

	return args
}
