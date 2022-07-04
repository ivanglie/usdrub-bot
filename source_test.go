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
	go s.UpdateRate()
	go s.GetRatef()
}

func TestUpdateRate(t *testing.T) {
	s.UpdateRate()
	if got := s.GetRate(); got < 0.0 {
		t.Errorf("Expected: greater or equal 0, got: %2f", got)
	}
}

func TestUpdateRateError(t *testing.T) {
	e.UpdateRate()
	if got := e.GetErr(); got == nil {
		t.Errorf("Expected: not nil, got: %v", got)
	}
}

func TestGetRatef(t *testing.T) {
	s.UpdateRate()
	got := s.GetRatef()
	want := fmt.Sprintf(s.GetPattern(), s.GetRate())
	if got != want {
		t.Errorf("Expected: %s, got: %s", want, got)
	}
}

func TestGetRatefError(t *testing.T) {
	e.UpdateRate()
	got := e.GetRatef()
	if e.err == nil || e.GetRate() > 0.0 {
		t.Errorf("Expected: error, got: %s", got)
	}
	want := fmt.Sprintf("%s error: %s", e.GetName(), e.GetErr().Error())
	if got != want {
		t.Errorf("Expected: error, got: %s", got)
	}
}
