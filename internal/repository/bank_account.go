package repository

import (
	"be-assessment-test/internal/types"
	"context"

	"github.com/go-errors/errors"
	"github.com/jmoiron/sqlx"
)

type BankAccount interface {
	CreateTx(ctx context.Context, tx *sqlx.Tx, req types.BankAccount) error
}

type bankAccountImpl struct {
	db *sqlx.DB
}

func NewBankAccount(db *sqlx.DB) BankAccount {
	return &bankAccountImpl{
		db: db,
	}
}

func (r *bankAccountImpl) CreateTx(ctx context.Context, tx *sqlx.Tx, req types.BankAccount) error {
	query := `
		INSERT INTO bank_accounts (
			id,
			user_id,
			created_at
		)
		VALUES (
			:id,
			:user_id,
			:created_at
		)
	`

	_, err := tx.NamedExecContext(ctx, query, req)
	if err != nil {
		return errors.New(err)
	}

	return nil
}
