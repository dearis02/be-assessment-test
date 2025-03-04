package handler

import (
	"be-assessment-test/internal/service"
	"be-assessment-test/internal/types"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Registration struct {
	registrationSvc service.Registration
}

func NewRegistration(registrationSvc service.Registration) *Registration {
	return &Registration{
		registrationSvc,
	}
}

func RegistrationRoutes(e *echo.Echo, registrationHandler *Registration) {
	e.POST("/v1/registration", registrationHandler.Register)
}

func (h *Registration) Register(c echo.Context) error {
	var req types.RegistrationCreateReq

	err := c.Bind(&req)
	if err != nil {
		return err
	}

	res, err := h.registrationSvc.Create(c.Request().Context(), req)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, types.ApiResponse{
		StatusCode: http.StatusCreated,
		Data:       res,
	})
}
