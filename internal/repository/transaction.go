package repository

import (
	"be-assessment-test/internal/types"
	"context"

	"github.com/go-errors/errors"
	"github.com/jmoiron/sqlx"
)

type Transaction interface {
	CreateTx(ctx context.Context, tx *sqlx.Tx, req types.Transaction) error
}

type transactionImpl struct {
	db *sqlx.DB
}

func NewTransaction(db *sqlx.DB) Transaction {
	return &transactionImpl{
		db: db,
	}
}

func (t *transactionImpl) CreateTx(ctx context.Context, tx *sqlx.Tx, req types.Transaction) error {
	query := `
		INSERT INTO transactions (
			id,
			bank_account_id,
			type,
			status,
			created_at
		)
		VALUES (
			:id,
			:bank_account_id,
			:type,
			:status,
			:created_at
		)
	`

	_, err := tx.NamedExecContext(ctx, query, req)
	if err != nil {
		return errors.New(err)
	}

	return nil
}
