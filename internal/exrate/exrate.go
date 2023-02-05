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
	er := &Rate{}
	er.Lock()
	defer er.Unlock()
	er.rate = rate
	er.err = err

	return er
}

// Get formatted exchange rate.
func (er *Rate) String() string {
	er.RLock()
	defer er.RUnlock()

	return fmt.Sprintf("%.2f RUB", er.rate)
}
