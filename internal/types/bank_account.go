package types

import (
	"net/http"
	"time"

	"github.com/go-errors/errors"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/volatiletech/null/v9"
)

// region repo types

type BankAccount struct {
	ID            uuid.UUID       `db:"id"`
	UserID        uuid.UUID       `db:"user_id"`
	AccountNumber string          `db:"account_number"`
	Balance       decimal.Decimal `db:"balance"`
	CreatedAt     time.Time       `db:"created_at"`
	UpdatedAt     null.Time       `db:"updated_at"`
}

// endregion repo types

// region service types

type BankAccountGetBalanceReq struct {
	AccountNumber string `param:"account-number"`
}

func (r BankAccountGetBalanceReq) Validate() error {
	err := validation.ValidateStruct(&r,
		validation.Field(&r.AccountNumber, validation.Required),
	)

	if err != nil {
		return errors.New(AppErr{
			StatusCode: http.StatusBadRequest,
			Message:    "account number param is required",
		})
	}

	return nil
}

type BankAccountGetBalanceRes struct {
	Balance decimal.Decimal `json:"balance"`
}

// endregion service types
