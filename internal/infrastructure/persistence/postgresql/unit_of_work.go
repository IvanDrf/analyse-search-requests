package postgresql

import (
	"context"
	"database/sql"

	"github.com/IvanDrf/analyse-search-requests/internal/domain/models"
)

const TRANSACTION = "transaction"

type unitOfWork struct {
	db *sql.DB
}

func NewUnitOfWork(db *sql.DB) *unitOfWork {
	return &unitOfWork{
		db: db,
	}
}

func (uof *unitOfWork) Execute(ctx context.Context, fn func(ctx context.Context) error) error {
	tx, err := uof.db.BeginTx(ctx, nil)
	if err != nil {
		return models.Error{
			Message: "can't start new transaction",
			Code:    models.ErrCodeInternal,
		}
	}

	ctx = context.WithValue(ctx, TRANSACTION, tx)

	if err := fn(ctx); err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}
