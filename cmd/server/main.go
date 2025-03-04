package main

import (
	"be-assessment-test/internal/config"
	"be-assessment-test/internal/types"
	"be-assessment-test/internal/util"
	"be-assessment-test/internal/util/dbutil"
	"be-assessment-test/internal/util/httputil"
	"context"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-errors/errors"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

func main() {
	cfg := config.New()
	util.InitLogger(cfg)

	db, err := dbutil.NewPostgres(&cfg.Database)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to database")
	}

	e := echo.New()
	e.HidePort = true
	e.HideBanner = true

	e.HTTPErrorHandler = httputil.HTTPErrorHandler(e)

	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, types.ApiResponse{
			StatusCode: http.StatusOK,
			Message:    "Hello",
		})
	})

	go func() {
		var err error

		address := cfg.Server.Address()
		e.Listener, err = net.Listen("tcp", address)
		if err != nil {
			log.Fatal().Err(err).Send()
		}

		log.Info().Msg("http server started on " + address)

		err = e.Start(address)
		if err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Send()
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit

	log.Info().Msg("shutting down the server")

	err = e.Shutdown(context.Background())
	if err != nil {
		log.Fatal().Err(err).Msg("failed to shutdown the server")
	}

	log.Info().Msg("closing db connections")

	err = db.Close()
	if err != nil {
		log.Error().Stack().Err(errors.New(err)).Msg("failed to close main db connection")
	}

	log.Info().Msg("db connections closed")

	log.Info().Msg("server shutdown successfully")
}
