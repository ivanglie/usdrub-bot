package utils

import (
	"testing"
)

func Test_setLogger(t *testing.T) {
	err := SetLogger(log)
	if err != nil {
		t.Error(err)
	}
}

func Test_setLogger_Error(t *testing.T) {
	err := SetLogger(nil)
	if err == nil {
		t.Error(err)
	}
}
