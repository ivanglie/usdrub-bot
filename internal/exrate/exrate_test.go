package exrate

import (
	"errors"
	"math/rand"
	"sync"
	"testing"
	"time"
)

func TestRate_dataRace(t *testing.T) {
	c := NewRate(func() (float64, error) { return 100 * rand.Float64(), nil })
	go func() {
		for {
			c.Update(nil)
		}
	}()

	for i := 0; i < 10; i++ {
		c.Rate()
		time.Sleep(100 * time.Millisecond)
	}
}

func TestRate_String(t *testing.T) {
	r := &Rate{}
	r.rate = 200.0

	if got := r.String(); len(got) == 0 {
		t.Errorf("Rate.String() = %v", got)
	}
}

func TestRate_Update(t *testing.T) {
	c := &Rate{rateFunc: func() (float64, error) { return 0, errors.New("error") }}
	c.Update(&sync.WaitGroup{})
}
