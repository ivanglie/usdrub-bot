package main

import (
	"errors"
	"fmt"
	"math/rand"
	"testing"
)

// Random Source
var s = Source{
	name:     "Random",
	pattern:  "1 US Dollar equals %.2f RUB by Random",
	rateFunc: func() (float64, error) { return 100 * rand.Float64(), nil },
}

// Error Source
var e = Source{
	name:     "Error source",
	pattern:  "1 US Dollar equals %.2f RUB by Error",
	rateFunc: func() (float64, error) { return 0.0, errors.New("Service error") },
}

// DATA RACE test
func TestDataRace(t *testing.T) {
	go s.updateRate()
	go s.getRatef()
}

func TestUpdateRate(t *testing.T) {
	s.updateRate()
	if got := s.rate; got < 0.0 {
		t.Errorf("Expected: greater or equal 0, got: %2f", got)
	}
}

func TestUpdateRateError(t *testing.T) {
	e.updateRate()
	if got := e.err; got == nil {
		t.Errorf("Expected: not nil, got: %v", got)
	}
}

func TestGetRatef(t *testing.T) {
	s.updateRate()
	got := s.getRatef()
	want := fmt.Sprintf(s.pattern, s.rate)
	if got != want {
		t.Errorf("Expected: %s, got: %s", want, got)
	}
}

func TestGetRatefError(t *testing.T) {
	e.updateRate()
	got := e.getRatef()
	if e.err == nil || e.rate > 0.0 {
		t.Errorf("Expected: error, got: %s", got)
	}
	want := fmt.Sprintf("%s error: %s", e.name, e.err.Error())
	if got != want {
		t.Errorf("Expected: error, got: %s", got)
	}
}
