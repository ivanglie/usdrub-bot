package scheduler

import (
	"testing"
)

func Test_setLogger(t *testing.T) {
	err := setLogger(log)
	if err != nil {
		t.Error(err)
	}
}

func Test_setLogger_Error(t *testing.T) {
	err := setLogger(nil)
	if err == nil {
		t.Error(err)
	}
}
