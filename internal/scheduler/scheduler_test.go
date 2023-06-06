package scheduler

import (
	"testing"

	"github.com/ivanglie/usdrub-bot/internal/logger"
)

func TestStartCmdOnSchedule(t *testing.T) {
	logger.Debug = true
	if err := StartCmdOnSchedule(func() {}); err != nil {
		t.Errorf("StartCmdOnSchedule() error = %v", err)
	}
}
