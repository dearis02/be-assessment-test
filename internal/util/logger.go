package util

import (
	"be-assessment-test/internal/config"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/go-errors/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var jsonIndentWriterBufPool = sync.Pool{
	New: func() interface{} {
		return bytes.NewBuffer(make([]byte, 0, 100))
	},
}

type jsonIndentWriter struct {
	Out io.Writer
}

func (w jsonIndentWriter) Write(p []byte) (int, error) {
	var buf = jsonIndentWriterBufPool.Get().(*bytes.Buffer)
	defer func() {
		buf.Reset()
		jsonIndentWriterBufPool.Put(buf)
	}()

	err := json.Indent(buf, p, "", "  ")
	if err != nil {
		return 0, err
	}

	_, err = buf.WriteTo(w.Out)
	return len(p), err
}

func InitLogger(c *config.Config) {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.ErrorStackMarshaler = ZeroLogErrorStackMarshaler
	zerolog.DefaultContextLogger = &log.Logger

	logCtx := log.With().Stack()

	if c.Environment != "production" {
		log.Logger = logCtx.Logger().Output(jsonIndentWriter{Out: os.Stderr})
		return
	}

	a := zerolog.Dict().Str("name", "be-assessment-test")
	log.Logger = logCtx.Dict("app", a).Logger()
}

func ZeroLogErrorStackMarshaler(err error) interface{} {
	var sterr *errors.Error
	var ok bool
	for err != nil {
		sterr, ok = err.(*errors.Error)
		if ok {
			break
		}

		u, ok := err.(interface {
			Unwrap() error
		})
		if !ok {
			return nil
		}

		err = u.Unwrap()
	}
	if sterr == nil {
		return nil
	}

	s := []string{}
	for _, frame := range sterr.StackFrames() {
		s = append(s, fmt.Sprintf("%s:%d", frame.File, frame.LineNumber))
	}
	return s
}
