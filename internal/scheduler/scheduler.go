package scheduler

import (
	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
)

// Start cmd on schedule
func StartCmdOnSchedule(cmd func(), spec string) {
	if spec == "" {
		spec = "0/5 * * * 1-5" // See https://crontab.guru/
	}
	log.Printf("Cron spec = %s", spec)

	scheduler := cron.New()
	defer scheduler.Stop()
	scheduler.AddFunc(spec, cmd)
	go scheduler.Start()
}
