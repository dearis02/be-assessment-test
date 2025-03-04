package types

import (
	"github.com/go-errors/errors"
)

type AppErr struct {
	StatusCode int
	Message    string
}

func (e AppErr) Error() string {
	return e.Message
}

var (
	ErrNoData = errors.New("no data")
)
