package handler

import (
	"be-assessment-test/internal/service"
	"be-assessment-test/internal/types"
	"net/http"

	"github.com/labstack/echo/v4"
)

type BankAccount struct {
	bankAccountSvc service.BankAccount
}

func NewBankAccount(bankAccountSvc service.BankAccount) *BankAccount {
	return &BankAccount{
		bankAccountSvc,
	}
}

func BankAccountRoutes(e *echo.Echo, bankAccountHandler *BankAccount) {
	e.GET("/v1/bank-accounts/:account-number", bankAccountHandler.GetBalance)
}

func (h BankAccount) GetBalance(c echo.Context) error {
	var req types.BankAccountGetBalanceReq

	err := c.Bind(&req)
	if err != nil {
		return err
	}

	res, err := h.bankAccountSvc.GetBalance(c.Request().Context(), req)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, types.ApiResponse{
		StatusCode: http.StatusOK,
		Data:       res,
	})
}
