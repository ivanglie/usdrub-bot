package scheduler

import (
	"os"

	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
)

var Debug bool

// Start cmd on schedule
func StartCmdOnSchedule(cmd func()) {
	spec := os.Getenv("CRON_SPEC")
	if spec == "" {
		spec = "0/5 * * * 1-5" // See https://crontab.guru/
	}
	if Debug {
		log.Printf("Cron spec = %s", spec)
	}

	scheduler := cron.New()
	defer scheduler.Stop()
	scheduler.AddFunc(spec, cmd)
	go scheduler.Start()
}
