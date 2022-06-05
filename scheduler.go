package main

import (
	"github.com/robfig/cron/v3"
)

// Start cmd on schedule
func startCmdOnSchedule(cmd func()) {
	spec := opts.CronSpec
	if spec == "" {
		spec = "0/5 * * * 1-5" // See https://crontab.guru/
	}
	log.Debugf("Cron spec = %s", spec)

	scheduler := cron.New()
	defer scheduler.Stop()
	scheduler.AddFunc(spec, cmd)
	go scheduler.Start()
}
