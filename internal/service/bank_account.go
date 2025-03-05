package service

import (
	"be-assessment-test/internal/repository"
	"be-assessment-test/internal/types"
	"context"
	"net/http"

	"github.com/go-errors/errors"
)

type BankAccount interface {
	GetBalance(ctx context.Context, req types.BankAccountGetBalanceReq) (types.BankAccountGetBalanceRes, error)
}

type bankAccountImpl struct {
	bankAccountRepo repository.BankAccount
}

func NewBankAccount(bankAccountRepo repository.BankAccount) BankAccount {
	return &bankAccountImpl{
		bankAccountRepo,
	}
}

func (s *bankAccountImpl) GetBalance(ctx context.Context, req types.BankAccountGetBalanceReq) (types.BankAccountGetBalanceRes, error) {
	res := types.BankAccountGetBalanceRes{}

	err := req.Validate()
	if err != nil {
		return res, err
	}

	bankAccount, err := s.bankAccountRepo.FindByBankAccountNumber(ctx, req.AccountNumber)
	if errors.Is(err, types.ErrNoData) {
		return res, errors.New(types.AppErr{StatusCode: http.StatusBadRequest, Message: "bank account not found"})
	} else if err != nil {
		return res, err
	}

	res.Balance = bankAccount.Balance

	return res, nil
}
