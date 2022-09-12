package scheduler

import (
	"os"

	"github.com/robfig/cron/v3"
)

var Debug bool

// Start cmd on schedule
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

	c.AddFunc(spec, cmd)
	go c.Start()
}
