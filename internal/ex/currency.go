package ex

import (
	"sync"
)

// Currency of exchange rate
type Currency struct {
	sync.RWMutex
	rate     float64
	err      error
	rateFunc func() (float64, error)
}

func New(rateFunc func() (float64, error)) *Currency {
	return &Currency{
		rateFunc: rateFunc,
	}
}

func (c *Currency) Update() {
	c.Lock()
	defer c.Unlock()
	c.rate, c.err = c.rateFunc()
}

// Get exchange rate
func (c *Currency) Rate() (float64, error) {
	c.RLock()
	defer c.RUnlock()
	return c.rate, c.err
}
