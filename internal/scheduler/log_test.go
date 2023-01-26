package scheduler

import (
	"testing"
)

func TestSetLogger(t *testing.T) {
	err := SetLogger(log)
	if err != nil {
		t.Error(err)
	}
}

func TestSetLogger_Error(t *testing.T) {
	err := SetLogger(nil)
	if err == nil {
		t.Error(err)
	}
}
