package exrate

import (
	"fmt"
	"sync"
)

type Rate struct {
	sync.RWMutex
	rate float64
	err  error
}

// New exchange rate.
func NewRate(rate float64, err error) *Rate {
	if rate == 0 || err != nil {
		return &Rate{}
	}

	er := &Rate{}
	er.Lock()
	defer er.Unlock()
	er.rate = rate
	er.err = err

	return er
}

// Update exchange rate.
func UpdateRate(f func() (float64, error)) *Rate {
	r, err := f()
	return NewRate(r, err)
}

// Get formatted exchange rate.
func (er *Rate) String() string {
	er.RLock()
	defer er.RUnlock()

	return fmt.Sprintf("%.2f RUB", er.rate)
}
