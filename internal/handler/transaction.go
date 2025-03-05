package handler

import (
	"be-assessment-test/internal/service"
	"be-assessment-test/internal/types"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Transaction struct {
	transactionSvc service.Transaction
}

func NewTransaction(transactionSvc service.Transaction) *Transaction {
	return &Transaction{
		transactionSvc,
	}
}

func TransactionRoutes(e *echo.Echo, transactionHandler *Transaction) {
	e.POST("/v1/transactions/_deposit", transactionHandler.Deposit)
	e.POST("/v1/transactions/_withdraw", transactionHandler.Withdraw)
}

func (t *Transaction) Deposit(c echo.Context) error {
	var req types.TransactionDepositReq

	err := c.Bind(&req)
	if err != nil {
		return err
	}

	res, err := t.transactionSvc.Deposit(c.Request().Context(), req)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, types.ApiResponse{
		StatusCode: http.StatusCreated,
		Data:       res,
	})
}

func (t *Transaction) Withdraw(c echo.Context) error {
	var req types.TransactionWithdrawReq

	err := c.Bind(&req)
	if err != nil {
		return err
	}

	res, err := t.transactionSvc.Withdraw(c.Request().Context(), req)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, types.ApiResponse{
		StatusCode: http.StatusCreated,
		Data:       res,
	})
}
