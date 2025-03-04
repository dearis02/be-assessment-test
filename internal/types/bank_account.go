package types

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/volatiletech/null/v9"
)

// region repo types

type BankAccount struct {
	ID        uuid.UUID       `db:"id"`
	UserID    uuid.UUID       `db:"user_id"`
	Balance   decimal.Decimal `db:"balance"`
	CreatedAt time.Time       `db:"created_at"`
	UpdatedAt null.Time       `db:"updated_at"`
}

// endregion repo types
