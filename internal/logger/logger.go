package logger

import (
	"errors"
	"io"
	stdlog "log"
	"os"
)

// Logger is an interface that represents the required methods to log data.
type Logger interface {
	Println(v ...interface{})
	Printf(format string, v ...interface{})
}

var (
	log   Logger = stdlog.New(io.MultiWriter(os.Stdout, os.Stderr), "", stdlog.LstdFlags)
	Debug bool
)

// SetLogger specifies the logger that the package should use.
func SetLogger(logger Logger) error {
	if logger == nil {
		return errors.New("logger is nil")
	}

	log = logger

	return nil
}
