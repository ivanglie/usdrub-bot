package exrate

import (
	"math/rand"
	"testing"
)

func TestNewRate(t *testing.T) {
	if got := NewRate(100*rand.Float64(), nil); got == nil {
		t.Errorf("NewRate() = %v", got)
	}
}

func TestUpdateRate(t *testing.T) {
	if got := UpdateRate(func() (float64, error) { return rand.Float64(), nil }); got.rate == 0 || got.err != nil {
		t.Errorf("UpdateRate() = %v", got)
	}
}

func TestRate_String(t *testing.T) {
	r := &Rate{}
	r.rate = 200.0

	if got := r.String(); len(got) == 0 {
		t.Errorf("Rate.String() = %v", got)
	}
}
