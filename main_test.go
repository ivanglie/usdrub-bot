package main

import (
	"testing"

	"github.com/sirupsen/logrus"
)

func TestSetupLogWithErrorLevel(t *testing.T) {
	setupLog(false)
	want := logrus.ErrorLevel
	if got := log.Level; got != want {
		t.Errorf("Expected: %v, got: %v", want, got)
	}
}

func TestSetupLogWithDebugLevel(t *testing.T) {
	setupLog(true)
	want := logrus.DebugLevel
	if got := log.Level; got != want {
		t.Errorf("Expected: %v, got: %v", want, got)
	}
}
