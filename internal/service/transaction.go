package service

import (
	"be-assessment-test/internal/repository"
	"be-assessment-test/internal/types"
	"be-assessment-test/internal/util/dbutil"
	"context"
	"net/http"
	"time"

	"github.com/go-errors/errors"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// check migration. for testing purposes
var internalBankAccountNumber = "00000000000000000000"

type Transaction interface {
	Deposit(ctx context.Context, req types.TransactionDepositReq) (types.TransactionDepositRes, error)
	Withdraw(ctx context.Context, req types.TransactionWithdrawReq) (types.TransactionWithdrawRes, error)
}

type transactionImpl struct {
	db              *sqlx.DB
	bankAccountRepo repository.BankAccount
	transactionRepo repository.Transaction
	ledgerEntryRepo repository.LedgerEntry
}

func NewTransaction(db *sqlx.DB, bankAccountRepo repository.BankAccount, transactionRepo repository.Transaction, ledgerEntryRepo repository.LedgerEntry) Transaction {
	return &transactionImpl{
		db,
		bankAccountRepo,
		transactionRepo,
		ledgerEntryRepo,
	}
}

func (s *transactionImpl) Deposit(ctx context.Context, req types.TransactionDepositReq) (types.TransactionDepositRes, error) {
	res := types.TransactionDepositRes{}

	err := req.Validate()
	if err != nil {
		return res, err
	}

	bankAccountExists := true
	userBankAccount, err := s.bankAccountRepo.FindByBankAccountNumber(ctx, req.BankAccountNumber)
	if errors.Is(err, types.ErrNoData) {
		bankAccountExists = false
	} else if err != nil {
		return res, err
	}

	if !bankAccountExists {
		return res, errors.New(types.AppErr{StatusCode: http.StatusBadRequest, Message: "bank account not found"})
	}

	internalBankAccount, err := s.bankAccountRepo.FindByBankAccountNumber(ctx, internalBankAccountNumber)
	if errors.Is(err, types.ErrNoData) {
		bankAccountExists = false
	} else if err != nil {
		return res, err
	}

	if !bankAccountExists {
		return res, errors.Errorf("internal bank account not found: account_number %s", internalBankAccountNumber)
	}

	timeNow := time.Now()

	id, err := uuid.NewV7()
	if err != nil {
		return res, errors.New(err)
	}

	trx := types.Transaction{
		ID:            id,
		BankAccountID: userBankAccount.ID,
		Type:          types.TransactionTypeDeposit,
		Status:        types.TransactionStatusSuccess,
		CreatedAt:     timeNow,
	}

	id, err = uuid.NewV7()
	if err != nil {
		return res, errors.New(err)
	}
	userLedger := types.LedgerEntry{
		ID:            id,
		BankAccountID: userBankAccount.ID,
		TransactionID: trx.ID,
		Type:          types.LedgerEntryTypeCredit,
		Amount:        req.Amount,
		BalanceAfter:  userBankAccount.Balance.Add(req.Amount),
		Description:   "deposit",
		CreatedAt:     timeNow,
	}

	id, err = uuid.NewV7()
	if err != nil {
		return res, errors.New(err)
	}
	internalLedger := types.LedgerEntry{
		ID:            id,
		BankAccountID: internalBankAccount.ID,
		TransactionID: trx.ID,
		Type:          types.LedgerEntryTypeDebit,
		Amount:        req.Amount,
		BalanceAfter:  internalBankAccount.Balance.Add(req.Amount),
		Description:   "deposit",
		CreatedAt:     timeNow,
	}

	tx, err := dbutil.NewSqlxTx(ctx, s.db, nil)
	if err != nil {
		return res, errors.New(err)
	}

	defer tx.Rollback()

	err = s.transactionRepo.CreateTx(ctx, tx, trx)
	if err != nil {
		return res, err
	}

	err = s.ledgerEntryRepo.CreateTx(ctx, tx, []types.LedgerEntry{userLedger, internalLedger})
	if err != nil {
		return res, err
	}

	userBankAccount.Balance = userBankAccount.Balance.Add(req.Amount)
	err = s.bankAccountRepo.UpdateBalanceTx(ctx, tx, userBankAccount)
	if err != nil {
		return res, err
	}

	internalBankAccount.Balance = internalBankAccount.Balance.Add(req.Amount)
	err = s.bankAccountRepo.UpdateBalanceTx(ctx, tx, internalBankAccount)
	if err != nil {
		return res, err
	}

	err = tx.Commit()
	if err != nil {
		return res, errors.New(err)
	}

	res.CurrentBalance = userBankAccount.Balance

	return res, nil
}

func (s *transactionImpl) Withdraw(ctx context.Context, req types.TransactionWithdrawReq) (types.TransactionWithdrawRes, error) {
	res := types.TransactionWithdrawRes{}

	err := req.Validate()
	if err != nil {
		return res, err
	}

	bankAccountExists := true
	userBankAccount, err := s.bankAccountRepo.FindByBankAccountNumber(ctx, req.BankAccountNumber)
	if errors.Is(err, types.ErrNoData) {
		bankAccountExists = false
	} else if err != nil {
		return res, err
	}

	if !bankAccountExists {
		return res, errors.New(types.AppErr{StatusCode: http.StatusBadRequest, Message: "bank account not found"})
	}
	if userBankAccount.Balance.LessThan(req.Amount) {
		return res, errors.New(types.AppErr{StatusCode: http.StatusBadRequest, Message: "insufficient balance"})
	}

	internalBankAccount, err := s.bankAccountRepo.FindByBankAccountNumber(ctx, internalBankAccountNumber)
	if errors.Is(err, types.ErrNoData) {
		bankAccountExists = false
	} else if err != nil {
		return res, err
	}

	if !bankAccountExists {
		return res, errors.Errorf("internal bank account not found: account_number %s", internalBankAccountNumber)
	}

	timeNow := time.Now()

	id, err := uuid.NewV7()
	if err != nil {
		return res, errors.New(err)
	}

	trx := types.Transaction{
		ID:            id,
		BankAccountID: userBankAccount.ID,
		Type:          types.TransactionTypeWithdraw,
		Status:        types.TransactionStatusSuccess,
		CreatedAt:     timeNow,
	}

	id, err = uuid.NewV7()
	if err != nil {
		return res, errors.New(err)
	}
	userLedger := types.LedgerEntry{
		ID:            id,
		BankAccountID: userBankAccount.ID,
		TransactionID: trx.ID,
		Type:          types.LedgerEntryTypeDebit,
		Amount:        req.Amount,
		BalanceAfter:  userBankAccount.Balance.Sub(req.Amount),
		Description:   "withdrawal",
		CreatedAt:     timeNow,
	}

	id, err = uuid.NewV7()
	if err != nil {
		return res, errors.New(err)
	}
	internalLedger := types.LedgerEntry{
		ID:            id,
		BankAccountID: internalBankAccount.ID,
		TransactionID: trx.ID,
		Type:          types.LedgerEntryTypeCredit,
		Amount:        req.Amount,
		BalanceAfter:  internalBankAccount.Balance.Sub(req.Amount),
		Description:   "withdrawal",
		CreatedAt:     timeNow,
	}

	tx, err := dbutil.NewSqlxTx(ctx, s.db, nil)
	if err != nil {
		return res, errors.New(err)
	}

	defer tx.Rollback()

	err = s.transactionRepo.CreateTx(ctx, tx, trx)
	if err != nil {
		return res, err
	}

	err = s.ledgerEntryRepo.CreateTx(ctx, tx, []types.LedgerEntry{userLedger, internalLedger})
	if err != nil {
		return res, err
	}

	userBankAccount.Balance = userBankAccount.Balance.Sub(req.Amount)
	err = s.bankAccountRepo.UpdateBalanceTx(ctx, tx, userBankAccount)
	if err != nil {
		return res, err
	}

	internalBankAccount.Balance = internalBankAccount.Balance.Sub(req.Amount)
	err = s.bankAccountRepo.UpdateBalanceTx(ctx, tx, internalBankAccount)
	if err != nil {
		return res, err
	}

	err = tx.Commit()
	if err != nil {
		return res, errors.New(err)
	}

	res.CurrentBalance = userBankAccount.Balance

	return res, nil
}
