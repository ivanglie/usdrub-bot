package scheduler

import (
	"os"

	"github.com/robfig/cron/v3"
)

var Debug bool

// StartCmdOnSchedule specified by cmd.
func StartCmdOnSchedule(cmd func()) {
	spec := os.Getenv("CRON_SPEC")
	if spec == "" {
		spec = "* * * * 1-5" // See https://crontab.guru/
	}

	if Debug {
		log.Printf("Cron spec = %s\n", spec)
	}

	c := cron.New()
	defer c.Stop()

	_, err := c.AddFunc(spec, cmd)
	if err != nil {
		log.Println(err)
	}
	go c.Start()
}
