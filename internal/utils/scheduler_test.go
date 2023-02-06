package utils

import "testing"

func TestStartCmdOnSchedule(t *testing.T) {
	Debug = true
	if err := StartCmdOnSchedule(func() {}); err != nil {
		t.Errorf("StartCmdOnSchedule(func() {}) error = %v", err)
	}
}
