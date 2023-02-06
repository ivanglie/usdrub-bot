package utils

import (
	"os"
	"time"

	"github.com/robfig/cron/v3"
)

// StartCmdOnSchedule specified by cmd.
func StartCmdOnSchedule(cmd func()) (err error) {
	spec := os.Getenv("CRON_SPEC")
	if spec == "" {
		spec = "* * * * 1-5" // See https://crontab.guru/
	}

	if Debug {
		log.Printf("Cron spec = %s\n", spec)
	}

	moscowTime, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		return
	}

	c := cron.New(cron.WithLocation(moscowTime))
	defer c.Stop()

	if _, err = c.AddFunc(spec, cmd); err != nil {
		return
	}

	go c.Start()

	return
}
