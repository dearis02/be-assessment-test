package repository

import (
	"be-assessment-test/internal/types"
	"context"
	"database/sql"

	"github.com/go-errors/errors"
	"github.com/jmoiron/sqlx"
)

type BankAccount interface {
	CreateTx(ctx context.Context, tx *sqlx.Tx, req types.BankAccount) error
	IsExistsByAccountNumber(ctx context.Context, accountNumber string) (bool, error)
	FindByBankAccountNumber(ctx context.Context, accountNumber string) (types.BankAccount, error)
	UpdateBalanceTx(ctx context.Context, tx *sqlx.Tx, req types.BankAccount) error
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
			account_number,
			created_at
		)
		VALUES (
			:id,
			:user_id,
			:account_number,
			:created_at
		)
	`

	_, err := tx.NamedExecContext(ctx, query, req)
	if err != nil {
		return errors.New(err)
	}

	return nil
}

func (r *bankAccountImpl) IsExistsByAccountNumber(ctx context.Context, accountNumber string) (bool, error) {
	res := true

	query := `
		SELECT EXISTS (
			SELECT 1
			FROM bank_accounts
			WHERE account_number = $1
		)
	`

	err := r.db.GetContext(ctx, &res, query, accountNumber)
	if err != nil {
		return res, err
	}

	return res, nil
}

func (r *bankAccountImpl) FindByBankAccountNumber(ctx context.Context, accountNumber string) (types.BankAccount, error) {
	res := types.BankAccount{}

	query := `
		SELECT
			id,
			user_id,
			account_number,
			balance,
			created_at,
			updated_at
		FROM bank_accounts
		WHERE account_number = $1
	`

	err := r.db.GetContext(ctx, &res, query, accountNumber)
	if errors.Is(err, sql.ErrNoRows) {
		return res, errors.New(types.ErrNoData)
	} else if err != nil {
		return res, err
	}

	return res, nil
}

func (r *bankAccountImpl) UpdateBalanceTx(ctx context.Context, tx *sqlx.Tx, req types.BankAccount) error {
	query := `
		UPDATE bank_accounts
		SET balance = $1
		WHERE id = $2
	`

	_, err := tx.ExecContext(ctx, query, req.Balance, req.ID)
	if err != nil {
		return errors.New(err)
	}

	return nil
}
