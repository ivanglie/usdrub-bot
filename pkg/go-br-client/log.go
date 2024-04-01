package br

import (
	"errors"
	stdlog "log"
	"os"
)

// Logger is an interface that represents the required methods to log data.
type Logger interface {
	Println(v ...interface{})
	Printf(format string, v ...interface{})
}

var log Logger = stdlog.New(os.Stderr, "", stdlog.LstdFlags)

// SetLogger specifies the logger that the package should use.
func SetLogger(logger Logger) error {
	if logger == nil {
		return errors.New("logger is nil")
	}
	log = logger
	return nil
}
