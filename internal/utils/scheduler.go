package utils

import (
	"os"
	"time"

	"github.com/robfig/cron/v3"
)

// StartCmdOnSchedule specified by cmd.
func StartCmdOnSchedule(cmd func()) {
	spec := os.Getenv("CRON_SPEC")
	if spec == "" {
		spec = "* * * * 1-5" // See https://crontab.guru/
	}

	if Debug {
		log.Printf("Cron spec = %s\n", spec)
	}

	moscowTime, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		log.Println(err)
		return
	}

	c := cron.New(cron.WithLocation(moscowTime))
	defer c.Stop()

	_, err = c.AddFunc(spec, cmd)
	if err != nil {
		log.Println(err)
		return
	}

	go c.Start()
}
