package repository

import (
	"be-assessment-test/internal/types"
	"context"
	"database/sql"

	"github.com/go-errors/errors"
	"github.com/jmoiron/sqlx"
)

type User interface {
	CreateTx(ctx context.Context, tx *sqlx.Tx, req types.User) error
	FindByNationalIDNumberOrPhoneNumber(ctx context.Context, nationalIDNumber, phoneNumber string) (types.User, error)
}

type userImpl struct {
	db *sqlx.DB
}

func NewUser(db *sqlx.DB) User {
	return &userImpl{
		db: db,
	}
}

func (r *userImpl) CreateTx(ctx context.Context, tx *sqlx.Tx, req types.User) error {
	query := `
		INSERT INTO users(
			id,
			name,
			national_identity_number,
			phone_number,
			created_at
		)
		VALUES (
			:id,
			:name,
			:national_identity_number,
			:phone_number,
			:created_at
		)
	`

	_, err := tx.NamedExecContext(ctx, query, req)
	if err != nil {
		return errors.New(err)
	}

	return nil
}

func (r *userImpl) FindByNationalIDNumberOrPhoneNumber(ctx context.Context, nationalIDNumber, phoneNumber string) (types.User, error) {
	res := types.User{}

	query := `
		SELECT
			id,
			name,
			national_identity_number,
			phone_number,
			created_at
		FROM users
		WHERE national_identity_number = $1
			OR phone_number = $2
	`

	err := r.db.GetContext(ctx, &res, query, nationalIDNumber, phoneNumber)
	if errors.Is(err, sql.ErrNoRows) {
		return res, errors.New(types.ErrNoData)
	} else if err != nil {
		return res, err
	}

	return res, nil
}
