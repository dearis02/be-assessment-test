package types

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// region repo types

type Transaction struct {
	ID            uuid.UUID         `db:"id"`
	BankAccountID uuid.UUID         `db:"bank_account_id"`
	Type          TransactionType   `db:"type"`
	Status        TransactionStatus `db:"status"`
	CreatedAt     time.Time         `db:"created_at"`
}

type TransactionType string

const (
	TransactionTypeWithdraw TransactionType = "withdraw"
	TransactionTypeDeposit  TransactionType = "deposit"
)

type TransactionStatus string

const (
	TransactionStatusPending TransactionStatus = "pending"
	TransactionStatusSuccess TransactionStatus = "success"
	TransactionStatusFailed  TransactionStatus = "failed"
)

// endregion repo types

// region service types

type TransactionDepositReq struct {
	BankAccountNumber string          `json:"bank_account_number"`
	Amount            decimal.Decimal `json:"amount"`
}

func (r TransactionDepositReq) Validate() error {
	err := validation.ValidateStruct(&r,
		validation.Field(&r.BankAccountNumber, validation.Required),
	)

	if err != nil {
		return err
	}

	ve := validation.Errors{}

	if r.Amount.IsZero() {
		ve["amount"] = validation.NewError("amount", "amount is required")
	}

	if r.Amount.IsNegative() {
		ve["amount"] = validation.NewError("amount", "amount must be positive")
	}

	if len(ve) > 0 {
		return ve
	}

	return nil
}

type TransactionDepositRes struct {
	CurrentBalance decimal.Decimal `json:"current_balance"`
}

type TransactionWithdrawReq struct {
	BankAccountNumber string          `json:"bank_account_number"`
	Amount            decimal.Decimal `json:"amount"`
}

func (r TransactionWithdrawReq) Validate() error {
	err := validation.ValidateStruct(&r,
		validation.Field(&r.BankAccountNumber, validation.Required),
	)

	if err != nil {
		return err
	}

	ve := validation.Errors{}

	if r.Amount.IsZero() {
		ve["amount"] = validation.NewError("amount", "amount is required")
	}

	if r.Amount.IsNegative() {
		ve["amount"] = validation.NewError("amount", "amount must be positive")
	}

	if len(ve) > 0 {
		return ve
	}

	return nil
}

type TransactionWithdrawRes TransactionDepositRes

// endregion service types
