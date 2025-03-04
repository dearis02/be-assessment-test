package dbutil

import (
	"be-assessment-test/internal/config"
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func NewPostgres(cfg *config.PostgresConfig) (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", cfg.ConString)
	if err != nil {
		return nil, err
	}

	db.SetMaxIdleConns(cfg.MaxIdleCons)
	db.SetMaxOpenConns(cfg.MaxOpenCons)

	return db, nil
}

func ClosePostgresConnection(db *sqlx.DB) error {
	return db.Close()
}

func NewSqlxTx(ctx context.Context, db *sqlx.DB, opts *sql.TxOptions) (*sqlx.Tx, error) {
	tx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}

	return tx, nil
}
