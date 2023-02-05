package scheduler

import (
	"errors"
	stdlog "log"
	"os"
	"time"

	"github.com/robfig/cron/v3"
)

// Logger is an interface that represents the required methods to log data.
type Logger interface {
	Println(v ...interface{})
	Printf(format string, v ...interface{})
}

var log Logger = stdlog.New(os.Stderr, "", stdlog.LstdFlags)

var Debug bool

// StartCmdOnSchedule specified by cmd.
func StartCmdOnSchedule(cmd func(), logger Logger) {
	spec := os.Getenv("CRON_SPEC")
	if spec == "" {
		spec = "* * * * 1-5" // See https://crontab.guru/
	}

	setLogger(logger)

	if Debug {
		log.Printf("Cron spec = %s\n", spec)
	}

	moscowTime, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		log.Println(err)
	}

	c := cron.New(cron.WithLocation(moscowTime))
	defer c.Stop()

	_, err = c.AddFunc(spec, cmd)
	if err != nil {
		log.Println(err)
	}

	go c.Start()
}

// SetLogger specifies the logger that the package should use.
func setLogger(logger Logger) error {
	if logger == nil {
		return errors.New("logger is nil")
	}
	log = logger
	return nil
}
