package repository

import (
	"be-assessment-test/internal/types"
	"context"

	"github.com/go-errors/errors"
	"github.com/jmoiron/sqlx"
)

type LedgerEntry interface {
	CreateTx(ctx context.Context, tx *sqlx.Tx, req []types.LedgerEntry) error
}

type ledgerEntryImpl struct {
	db *sqlx.DB
}

func NewLedgerEntry(db *sqlx.DB) LedgerEntry {
	return &ledgerEntryImpl{
		db: db,
	}
}

func (l *ledgerEntryImpl) CreateTx(ctx context.Context, tx *sqlx.Tx, req []types.LedgerEntry) error {
	query := `
		INSERT INTO ledger_entries (
			id,
			bank_account_id,
			transaction_id,
			type,
			amount,
			balance_after,
			description,
			created_at
		)
		VALUES (
			:id,
			:bank_account_id,
			:transaction_id,
			:type,
			:amount,
			:balance_after,
			:description,
			:created_at
		)
	`

	_, err := tx.NamedExecContext(ctx, query, req)
	if err != nil {
		return errors.New(err)
	}

	return nil
}
