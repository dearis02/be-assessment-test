package types

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// region repo types

type LedgerEntry struct {
	ID            uuid.UUID       `db:"id"`
	BankAccountID uuid.UUID       `db:"bank_account_id"`
	TransactionID uuid.UUID       `db:"transaction_id"`
	Type          LedgerEntryType `db:"type"`
	Amount        decimal.Decimal `db:"amount"`
	BalanceAfter  decimal.Decimal `db:"balance_after"`
	Description   string          `db:"description"`
	CreatedAt     time.Time       `db:"created_at"`
}

type LedgerEntryType string

const (
	LedgerEntryTypeDebit  LedgerEntryType = "debit"
	LedgerEntryTypeCredit LedgerEntryType = "credit"
)

// endregion
