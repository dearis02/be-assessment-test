package httputil

import (
	"be-assessment-test/internal/types"
	"net/http"

	"github.com/go-errors/errors"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

func HTTPErrorHandler(e *echo.Echo) func(err error, c echo.Context) {
	return func(err error, c echo.Context) {
		if c.Response().Committed {
			return
		}

		code := http.StatusInternalServerError
		body := types.ApiResponse{
			Message: http.StatusText(code),
		}

		switch err := err.(type) {
		case *echo.HTTPError:
			if err.Internal != nil {
				if herr, ok := err.Internal.(*echo.HTTPError); ok {
					err = herr
				}
			}

			code = err.Code
			switch m := err.Message.(type) {
			case string:
				body.Message = m
				if e.Debug {
					body.Errors = map[string]string{
						"debug": err.Error(),
					}
				}
			case error:
				body.Message = m.Error()
			default:
				log.Warn().Err(err).Type("type", m).Msg("unhandled echo http error message type")
			}
		case *errors.Error:
			if d, ok := err.Err.(types.AppErr); ok {
				code = d.StatusCode
				body.Message = d.Message

				if d.Message == "" {
					body.Message = http.StatusText(code)
				}
			} else {
				log.Error().Stack().Err(err).Send()
			}
		case validation.Errors:
			code = http.StatusUnprocessableEntity
			errors := map[string]any{}
			for key, value := range err {
				errors[key] = []string{value.Error()}
			}
			body.Message = "Validation Error"
			body.Errors = errors
		}

		body.StatusCode = code

		if c.Request().Method == http.MethodHead {
			err = c.NoContent(code)
		} else {
			err = c.JSON(code, body)
		}
		if err != nil {
			log.Error().Err(err).Send()
		}
	}
}
